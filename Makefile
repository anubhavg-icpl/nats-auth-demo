.PHONY: help build run clean install docker-up docker-down test podman-up podman-down

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=nats-demo
GO=go
DOCKER_COMPOSE=docker-compose
PODMAN_COMPOSE=podman-compose

# Detect container runtime
CONTAINER_RUNTIME := $(shell command -v podman 2> /dev/null)
ifdef CONTAINER_RUNTIME
    COMPOSE_CMD=podman-compose -f podman-compose.yml
    RUNTIME_NAME=Podman
else
    COMPOSE_CMD=docker-compose
    RUNTIME_NAME=Docker
endif

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install Go dependencies
	$(GO) mod download
	$(GO) mod tidy

build: ## Build the demo application
	$(GO) build -o $(BINARY_NAME) cmd/main.go
	@echo "Built $(BINARY_NAME)"

run: build ## Build and run the demo application
	./$(BINARY_NAME)

clean: ## Remove built binaries
	rm -f $(BINARY_NAME)
	@echo "Cleaned build artifacts"

test: ## Run tests (if any)
	$(GO) test -v ./...

fmt: ## Format Go code
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

up: ## Start all NATS servers (auto-detects docker/podman)
	@echo "Starting NATS servers using $(RUNTIME_NAME)..."
	$(COMPOSE_CMD) up -d
	@echo ""
	@echo "✓ All NATS servers started with $(RUNTIME_NAME)"
	@echo "  - Basic Auth:       localhost:4222 (monitor: :8222)"
	@echo "  - Allow/Deny:       localhost:4223 (monitor: :8223)"
	@echo "  - Allow Responses:  localhost:4224 (monitor: :8224)"
	@echo "  - Queue Perms:      localhost:4225 (monitor: :8225)"
	@echo "  - Accounts:         localhost:4226 (monitor: :8226)"

down: ## Stop all NATS servers (auto-detects docker/podman)
	@echo "Stopping NATS servers using $(RUNTIME_NAME)..."
	$(COMPOSE_CMD) down
	@echo "✓ All NATS servers stopped"

logs: ## Show container logs (auto-detects docker/podman)
	$(COMPOSE_CMD) logs -f

ps: ## Show running containers (auto-detects docker/podman)
	$(COMPOSE_CMD) ps

# Docker-specific targets
docker-up: ## Start all NATS servers with docker-compose
	docker-compose up -d
	@echo "✓ All NATS servers started with Docker"

docker-down: ## Stop all NATS servers with docker-compose
	docker-compose down
	@echo "✓ All NATS servers stopped"

docker-logs: ## Show docker logs
	docker-compose logs -f

docker-basic: ## Start only basic auth server with docker
	docker-compose up -d nats-basic

docker-accounts: ## Start only accounts server with docker
	docker-compose up -d nats-accounts

# Podman-specific targets
podman-up: ## Start all NATS servers with podman-compose
	podman-compose -f podman-compose.yml up -d
	@echo "✓ All NATS servers started with Podman"

podman-down: ## Stop all NATS servers with podman-compose
	podman-compose -f podman-compose.yml down
	@echo "✓ All NATS servers stopped"

podman-logs: ## Show podman logs
	podman-compose -f podman-compose.yml logs -f

podman-basic: ## Start only basic auth server with podman
	podman-compose -f podman-compose.yml up -d nats-basic

podman-accounts: ## Start only accounts server with podman
	podman-compose -f podman-compose.yml up -d nats-accounts

podman-pull: ## Pull NATS images for podman
	podman pull docker.io/library/nats:latest

# Individual server targets using nats-server directly
start-basic: ## Start basic auth NATS server (requires nats-server)
	nats-server -c config/basic-auth.conf

start-allow-deny: ## Start allow/deny NATS server
	nats-server -c config/allow-deny.conf

start-allow-responses: ## Start allow-responses NATS server
	nats-server -c config/allow-responses.conf

start-queue-perms: ## Start queue permissions NATS server
	nats-server -c config/queue-permissions.conf

start-accounts: ## Start accounts NATS server
	nats-server -c config/accounts.conf

all: clean install build ## Clean, install dependencies, and build

.PHONY: start-basic start-allow-deny start-allow-responses start-queue-perms start-accounts
.PHONY: up down logs ps podman-up podman-down podman-logs podman-basic podman-accounts podman-pull
