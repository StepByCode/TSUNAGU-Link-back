.PHONY: help build run dev test clean docker-up docker-down docker-build migrate openapi-gen install-tools

APP_NAME=server
BIN_DIR=bin
CMD_DIR=cmd/server

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

build: ## Build the application
	@echo "Building application..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)

run: build ## Build and run the application locally
	@echo "Running application..."
	./$(BIN_DIR)/$(APP_NAME)

dev: ## Run the application in development mode with auto-reload (requires air)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)
	rm -f coverage.out coverage.html

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-build: ## Build Docker image
	docker-compose build

docker-logs: ## Show Docker logs
	docker-compose logs -f

migrate-up: ## Run database migrations up
	migrate -path db/migrations -database "postgresql://tsunagu:tsunagu_password@localhost:5432/tsunagu_db?sslmode=disable" up

migrate-down: ## Run database migrations down
	migrate -path db/migrations -database "postgresql://tsunagu:tsunagu_password@localhost:5432/tsunagu_db?sslmode=disable" down

migrate-create: ## Create a new migration file (usage: make migrate-create name=create_users_table)
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	migrate create -ext sql -dir db/migrations -seq $(name)

openapi-gen: ## Generate code from OpenAPI specification
	@echo "Generating code from OpenAPI spec..."
	oapi-codegen -package api -generate types,server,spec api/openapi.yaml > internal/handler/api_generated.go

deps: ## Download dependencies
	go mod download
	go mod tidy

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

lint: ## Run linter (requires golangci-lint)
	@which golangci-lint > /dev/null || (echo "Please install golangci-lint: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

setup: install-tools deps docker-up migrate-up ## Setup project (install tools, dependencies, start Docker, run migrations)
	@echo "Project setup complete!"

all: clean build test ## Clean, build and test
