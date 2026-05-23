package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PLT875/urlshortener/handlers"
	"github.com/PLT875/urlshortener/persistence"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/urlshortener?sslmode=disable"
	}

	repo, err := persistence.NewPostgresRepository(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()

	handler := handlers.NewHandler(repo)

	http.HandleFunc("/shorten", handler.ShortenURL)
	http.HandleFunc("/", handler.Redirect)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
