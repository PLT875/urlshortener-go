package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/PLT875/urlshortener/handlers"
	"github.com/PLT875/urlshortener/persistence"
)

func TestShortenerWithPostgresContainer(t *testing.T) {
	ctx := context.Background()

	// Start PostgreSQL container with alpine image
	postgresContainer, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "urlshortener",
				"POSTGRES_USER":     "user",
				"POSTGRES_PASSWORD": "password",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(10 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	defer postgresContainer.Terminate(ctx)

	// Get container connection details
	pgHost, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get postgres host: %v", err)
	}
	pgPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("failed to get postgres port: %v", err)
	}

	dbURL := fmt.Sprintf("postgres://user:password@%s:%s/urlshortener?sslmode=disable",
		pgHost, pgPort.Port())

	// Initialize repository with container database
	repo, err := persistence.NewPostgresRepository(dbURL)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer repo.Close()

	// Create handler with test repository
	handler := handlers.NewHandler(repo)

	// Create test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/shorten" {
			handler.ShortenURL(w, r)
		} else {
			handler.Redirect(w, r)
		}
	}))
	defer server.Close()

	// Test 1: Shorten URL
	t.Run("shorten URL", func(t *testing.T) {
		payload := []byte(`{"url":"https://example.com"}`)
		resp, err := http.Post(server.URL+"/shorten", "application/json", bytes.NewReader(payload))
		if err != nil {
			t.Fatalf("shorten request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
		}

		var result map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if result["short"] == "" {
			t.Fatal("short code is empty")
		}
	})

	// Test 2: Redirect to shortened URL
	t.Run("redirect to shortened URL", func(t *testing.T) {
		// First, create a short URL
		payload := []byte(`{"url":"https://github.com"}`)
		resp, err := http.Post(server.URL+"/shorten", "application/json", bytes.NewReader(payload))
		if err != nil {
			t.Fatalf("shorten request failed: %v", err)
		}

		var result map[string]string
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		shortCode := result["short"]

		// Now test redirect
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // Don't follow redirects
			},
		}

		resp, err = client.Get(server.URL + "/" + shortCode)
		if err != nil {
			t.Fatalf("redirect request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("unexpected status: %d, expected redirect", resp.StatusCode)
		}

		location := resp.Header.Get("Location")
		if location != "https://github.com" {
			t.Fatalf("unexpected redirect location: %s", location)
		}
	})
}
