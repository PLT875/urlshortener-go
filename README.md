# Go URL Shortener

This is a minimal URL shortener service written in Go with a TypeScript frontend, designed to use as few external libraries as possible. It provides endpoints to shorten URLs and redirect to the original URLs using PostgreSQL for persistence. Testing is performed using test containers for integration tests.

## Features
- Shorten URLs via HTTP API
- TypeScript frontend with form-based URL shortening
- Redirect shortened URLs to original URLs
- PostgreSQL persistence
- Minimal dependencies (Go: only PostgreSQL driver + testcontainers; Frontend: TypeScript only)
- Integration tests using test containers
- Makefile with convenient build/test/run targets

## Usage

### Prerequisites
- Go 1.21+
- Node.js 18+ (for frontend)
- PostgreSQL running locally or via Docker
- Docker & Docker Compose (for database containers)

### Quick Start with Make

```sh
# Backend: Build and run
make build
make run

# Frontend: Build TypeScript
make site-install   # Install npm dependencies (first time only)
make site-build     # Build TypeScript to dist/

# Database
make db-up          # Start PostgreSQL with docker-compose
make db-down        # Stop PostgreSQL
```

See all available commands:
```sh
make help
```

### Manual Setup

#### Backend

1. Start PostgreSQL:
   ```sh
   docker-compose up -d
   ```

2. Run the server:
   ```sh
   go run ./cmd/urlshortener/main.go
   ```
   
   Or set a custom database URL:
   ```sh
   DATABASE_URL="postgres://user:password@localhost:5432/urlshortener?sslmode=disable" go run ./cmd/urlshortener/main.go
   ```

3. Server runs at `http://localhost:8080`

#### Frontend

1. Install dependencies (first time):
   ```sh
   npm install
   ```

2. Build TypeScript:
   ```sh
   npm run build:site
   ```

3. Watch mode (auto-rebuild on changes):
   ```sh
   npm run dev:site
   ```

#### Using the Application

1. Open `http://localhost:8080` in your browser
2. Enter a URL in the form
3. Get a shortened URL
4. Or use the API:
   ```sh
   curl -X POST -H "Content-Type: application/json" \
     -d '{"url": "https://example.com/very/long/path"}' \
     http://localhost:8080/shorten
   ```

## Testing

### Go Backend

Integration tests use test containers (requires Docker). Run tests with:
```sh
make test
```

Or manually:
```sh
go test ./...
```

## Available Make Targets

### Backend
- `make test` - Run Go tests with Docker configuration
- `make build` - Build the application
- `make run` - Run the application (requires `make build` first)
- `make dev` - Build and run in one command
- `make clean` - Remove build artifacts
- `make tidy` - Tidy Go modules
- `make fmt` - Format Go code
- `make lint` - Run golangci-lint
- `make watch-test` - Watch for changes and auto-run tests

### Frontend (TypeScript/npm)
- `make site-install` - Install npm dependencies
- `make site-build` - Build TypeScript to dist/
- `make site-dev` - Watch TypeScript for changes
- `make site-clean` - Clean compiled output

### Database
- `make db-up` - Start PostgreSQL with docker-compose
- `make db-down` - Stop PostgreSQL
- `make db-logs` - Show PostgreSQL logs

## Integration Testing

Integration tests use test containers (requires Docker). Run tests with:
```sh
go test ./...
```

## Project Structure

```
├── Makefile                 # Build automation targets
├── docker-compose.yml       # PostgreSQL container setup
├── go.mod                   # Go dependencies
├── package.json             # npm/TypeScript dependencies
├── tsconfig.json            # TypeScript configuration
│
├── cmd/
│   └── urlshortener/
│       ├── main.go          # Entry point and HTTP server setup
│       └── main_test.go     # Integration tests with test containers
│
├── internal/
│   ├── domain/              # Core business logic
│   │   ├── shortener.go     # URL shortening algorithm
│   │   └── repository.go    # Storage interface
│   ├── handlers/            # HTTP handlers
│   │   └── handler.go       # Request/response handlers
│   └── persistence/         # Storage implementations
│       ├── postgres_repository.go  # PostgreSQL implementation
│       └── repository.go    # Repository interface
│
├── site/                    # TypeScript frontend
│   ├── src/
│   │   └── app.ts           # Frontend application logic
│   ├── public/
│   │   └── index.html       # HTML served to browser
│   └── dist/                # Compiled JavaScript (generated)
│
└── bin/
    └── urlshortener         # Compiled binary (generated)
```

## Environment Variables
- `DATABASE_URL`: PostgreSQL connection string (default: `postgres://user:password@localhost:5432/urlshortener?sslmode=disable`)

---

Easily extend with additional repository implementations without modifying domain or handler layers.

