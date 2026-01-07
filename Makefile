.PHONY: help setup start stop restart clean build run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Install Go dependencies
	@echo "📦 Installing Go dependencies..."
	go mod download
	go mod tidy

build: ## Build the Go example
	@echo "🔨 Building Go example..."
	go build -o bin/example main.go
	@echo "✅ Build successful!"

start: ## Start the Cribl server using Docker Compose
	@echo "🚀 Starting Cribl server..."
	docker-compose up -d
	@echo "✅ Cribl server started on http://localhost:9000"
	@echo "⏳ Waiting for server to be ready..."
	@sleep 10

stop: ## Stop the Cribl server
	@echo "🛑 Stopping Cribl server..."
	docker-compose down

restart: stop start ## Restart the Cribl server

clean: ## Stop and remove containers, volumes, and networks
	@echo "🧹 Cleaning up Docker resources..."
	docker-compose down -v
	@echo "🧹 Removing build artifacts..."
	@rm -rf bin/

run: build ## Run the Go example (builds first)
	@echo "🏃 Running Go example..."
	./bin/example

dev: start run ## Start server and run the example

