.PHONY: help build test run clean docker-build docker-up docker-down lint fmt

# Variables
APP_NAME=kalshi-api
WORKER_NAME=kalshi-worker
VERSION?=latest
GO=go
DOCKER=docker
DOCKER_COMPOSE=docker-compose

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build API and Worker binaries
	@echo "Building $(APP_NAME)..."
	$(GO) build -o bin/api ./cmd/api
	@echo "Building $(WORKER_NAME)..."
	$(GO) build -o bin/worker ./cmd/worker

test: ## Run tests
	@echo "Running tests..."
	$(GO) test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests and generate coverage report
	@echo "Generating coverage report..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

run-api: ## Run API server locally
	@echo "Starting API server..."
	$(GO) run ./cmd/api

run-worker: ## Run worker locally
	@echo "Starting worker..."
	$(GO) run ./cmd/worker

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	$(GO) fmt ./...
	gofmt -s -w .

tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	$(GO) mod tidy

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down

docker-logs: ## Show Docker logs
	$(DOCKER_COMPOSE) logs -f

docker-clean: docker-down ## Clean Docker resources
	@echo "Cleaning Docker resources..."
	$(DOCKER_COMPOSE) down -v
	$(DOCKER) system prune -f

deps: ## Install dependencies
	@echo "Installing dependencies..."
	$(GO) mod download

dev: ## Start development environment
	@echo "Starting development environment..."
	$(DOCKER_COMPOSE) up -d redis
	@echo "Redis is running on localhost:6379"
	@echo "Run 'make run-api' and 'make run-worker' in separate terminals"

all: clean deps build test ## Clean, install deps, build, and test

.DEFAULT_GOAL := help
