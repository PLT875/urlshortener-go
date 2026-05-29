.PHONY: help test build run clean tidy fmt lint db-up db-down db-logs dev watch-test

# Use bash shell (Git for Windows)
SHELL := bash
.SHELLFLAGS := -c

# Environment variables for testcontainers on Windows
export TESTCONTAINERS_RYUK_DISABLED := true
export DOCKER_HOST := npipe:////./pipe/docker_engine

help:
	@echo "Available targets:"
	@echo "  make test      - Run Go tests with proper Docker configuration"
	@echo "  make build     - Build the application"
	@echo "  make run       - Run the application"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make tidy      - Tidy Go modules"
	@echo "  make fmt       - Format Go code"
	@echo "  make lint      - Run golangci-lint (if installed)"
	@echo "  make db-up     - Start PostgreSQL with docker-compose"
	@echo "  make db-down   - Stop PostgreSQL"
	@echo "  make dev       - Build and run the application"

test:
	@echo "Running tests with Windows Docker configuration..."
	go test ./... -v -timeout 120s

build:
	@echo "Building application..."
	go build -o bin/urlshortener ./cmd/urlshortener

run: build
	@echo "Running application..."
	.bin/urlshortener

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
