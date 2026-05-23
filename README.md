# Go URL Shortener

This is a minimal URL shortener service written in Go, designed to use as few external libraries as possible. It provides endpoints to shorten URLs and redirect to the original URLs using PostgreSQL for persistence. Testing is performed using test containers for integration tests.

## Features
- Shorten URLs via HTTP API
- Redirect shortened URLs to original URLs
- PostgreSQL persistence
- Minimal dependencies (only PostgreSQL driver + testcontainers for testing)
- Integration tests using test containers

## Usage

### Prerequisites
- Go 1.21+
- PostgreSQL running locally or via environment variable

### Setup

1. Start PostgreSQL:
   ```sh
   docker run --name postgres -e POSTGRES_DB=urlshortener -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -p 5432:5432 postgres:15
   ```

2. Run the server:
   ```sh
   go run main.go
   ```
   
   Or set a custom database URL:
   ```sh
   DATABASE_URL="postgres://user:password@localhost:5432/urlshortener?sslmode=disable" go run main.go
   ```

3. Shorten a URL:
   ```sh
   curl -X POST -H "Content-Type: application/json" -d '{"url": "https://example.com"}' http://localhost:8080/shorten
   ```

4. Redirect:
   Visit `http://localhost:8080/<shortcode>` in your browser.

## Testing

Integration tests use test containers (requires Docker). Run tests with:
```sh
go test ./...
```

## Project Structure
- `main.go`: Entry point and HTTP server setup
- `domain/`: Core business logic
  - `shortener.go`: URL shortening algorithm
  - `repository.go`: Storage interface
- `handlers/`: HTTP handlers
  - `handler.go`: Request/response handlers
- `persistence/`: Storage implementations
  - `postgres_repository.go`: PostgreSQL implementation
  - `memory_repository.go`: In-memory reference implementation
- `main_test.go`: Integration tests with test containers

## Environment Variables
- `DATABASE_URL`: PostgreSQL connection string (default: `postgres://user:password@localhost:5432/urlshortener?sslmode=disable`)

---

Easily extend with additional repository implementations without modifying domain or handler layers.

