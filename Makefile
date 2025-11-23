# Makefile

# Variables
APP_NAME := gin-crud-demo
DOCKER_IMAGE := $(APP_NAME):latest
GOOSE_DRIVER := postgres
GOOSE_DBSTRING := "user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable host=$(DB_HOST) port=$(DB_PORT)"
GOOSE_MIGRATION_DIR := internal/database/migrations

.PHONY: all build run test test-cover mocks migrate-up migrate-down migrate-create docker-build docker-up docker-down help

all: build

build: ## Build the application binary
	go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

run: ## Run the application locally
	go run cmd/$(APP_NAME)/main.go

test: ## Run tests
	go test -v ./...

test-cover: ## Run tests with coverage report
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

mocks: ## Generate mocks using mockery
	mockery

migrate-up: ## Run database migrations up
	goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

migrate-down: ## Run database migrations down
	goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down

migrate-create: ## Create a new database migration
	@read -p "Enter migration name: " name; \
	goose -dir $(GOOSE_MIGRATION_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING) create $$name sql

docker-build: ## Build the Docker image
	docker build -t $(DOCKER_IMAGE) -f deployments/Dockerfile .

docker-up: ## Start Docker Compose services
	docker-compose -f deployments/docker-compose.yml up -d

docker-down: ## Stop Docker Compose services
	docker-compose -f deployments/docker-compose.yml down

help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
