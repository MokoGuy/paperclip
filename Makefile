.PHONY: build run generate clean test coverage help deps fmt dev-setup

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Build the paperclip binary
	@echo "Building paperclip..."
	@go build -o bin/paperclip ./cmd/paperclip

run: build ## Build and run the application
	@./bin/paperclip

generate: ## Generate SQLC database code
	@echo "Generating SQLC code..."
	@sqlc generate

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -cover ./...

deps: ## Install/update dependencies
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

dev-setup: deps generate ## Set up development environment
	@echo "Development environment ready!"
