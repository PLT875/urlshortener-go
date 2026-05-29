package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/PLT875/urlshortener/internal/domain"
)

type Handler struct {
	repo domain.Repository
}

func NewHandler(repo domain.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		URL string `json:"url"`
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &req); err != nil || !strings.HasPrefix(req.URL, "http") {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	code := domain.Shorten(req.URL)
	h.repo.Save(code, req.URL)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"short": code})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" || code == "shorten" {
		http.NotFound(w, r)
		return
	}
	url, ok := h.repo.Get(code)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
