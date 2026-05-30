.PHONY: help test build run clean tidy fmt lint db-up db-down db-logs dev watch-test \
	site-install site-build site-dev site-clean

# Use bash shell (Git for Windows)
SHELL := bash
.SHELLFLAGS := -c

help:
	@echo "Available targets:"
	@echo ""
	@echo "Go Backend:"
	@echo "  make test      - Run Go tests with proper Docker configuration"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make tidy      - Tidy Go modules"
	@echo "  make fmt       - Format Go code"
	@echo "  make lint      - Run golangci-lint (if installed)"
	@echo "  make dev       - Build and run the application"
	@echo "  make watch-test - Watch for changes and run tests"
	@echo ""
	@echo "TypeScript Frontend (site):"
	@echo "  make site-install    - Install npm dependencies"
	@echo "  make site-build      - Build TypeScript site"
	@echo "  make site-dev        - Watch and rebuild site on changes"
	@echo "  make site-clean      - Clean compiled site"
	@echo ""
	@echo "Database:"
	@echo "  make db-up     - Start PostgreSQL with docker-compose"
	@echo "  make db-down   - Stop PostgreSQL"
	@echo "  make db-logs   - Show database logs"

testwindows:
	@echo "Running tests with Windows Docker configuration..."
	# Environment variables for testcontainers on Windows
	TESTCONTAINERS_RYUK_DISABLED=true DOCKER_HOST=npipe:////./pipe/docker_engine go test ./... -v -timeout 120s

test:
	@echo "Running tests..."
	go test ./... -v -timeout 120s

build:
	@echo "Building application..."
	go build -o bin/urlshortener ./cmd/urlshortener

run:
	@echo "Running application..."
	./bin/urlshortener

clean:
	@echo "Cleaning build artifacts..."
	rm -f bin/*
	go clean

tidy:
	@echo "Tidying Go modules..."
	go mod tidy

fmt:
	@echo "Formatting Go code..."
	go fmt ./...

lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Database targets (for use with docker-compose)
db-up:
	@echo "Starting PostgreSQL container with docker-compose..."
	docker-compose up -d

db-down:
	@echo "Stopping PostgreSQL container..."
	docker-compose down

db-logs:
	@echo "Showing database logs..."
	docker-compose logs postgres

# Development targets
dev: build run

watch-test:
	@echo "Watching for changes and running tests..."
	@echo "Install 'reflex' first: go install github.com/cespare/reflex@latest"
	reflex -r '\.go$$' -s 'make test'

# Frontend targets (npm/TypeScript/site)
site-install:
	@echo "Installing npm dependencies..."
	npm install

site-build:
	@echo "Building TypeScript site..."
	npm run build:site

site-dev:
	@echo "Watching site for changes..."
	npm run dev:site

site-clean:
	@echo "Cleaning compiled site..."
	npm run clean:site
