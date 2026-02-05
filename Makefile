.PHONY: help install docker-up docker-down server web mobile clean db-migrate-up db-migrate-down db-migrate-version

# Auto-detect docker compose command (support both v1 and v2)
DOCKER_COMPOSE := $(shell command -v docker-compose 2>/dev/null || echo "docker compose")

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

install: ## Install all dependencies
	@echo "Installing root dependencies with pnpm..."
	pnpm install
	@echo "Installing Go dependencies..."
	cd apps/server && go mod tidy && go mod download
	@echo "Installing Flutter dependencies..."
	cd apps/mobile && flutter pub get
	@echo "All dependencies installed!"

docker-up: ## Start Docker containers (PostgreSQL + Redis)
	@echo "Starting Docker containers..."
	$(DOCKER_COMPOSE) up -d
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Applying database migrations..."
	@$(MAKE) db-migrate-up
	@echo "Docker containers are running!"

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	$(DOCKER_COMPOSE) down
	@echo "Docker containers stopped!"

db-migrate-up: ## Apply database migrations
	@echo "Applying database migrations..."
	cd apps/server && go run cmd/migrate/main.go up

db-migrate-down: ## Roll back last database migration
	@echo "Rolling back last database migration..."
	cd apps/server && go run cmd/migrate/main.go down

db-migrate-version: ## Show current migration version
	@echo "Current migration version..."
	cd apps/server && go run cmd/migrate/main.go version

docker-logs: ## View Docker container logs
	$(DOCKER_COMPOSE) logs -f

db-status: ## Check database status
	$(DOCKER_COMPOSE) exec db psql -U admin -d price_comparison -c "SELECT COUNT(*) FROM stores;"

server: ## Run Go backend server
	@echo "Starting Go server..."
	cd apps/server && go run cmd/main.go

web: ## Run Next.js web app with Turbopack
	@echo "Starting Next.js web app with Turbopack..."
	cd apps/web && pnpm run dev

mobile: ## Run Flutter mobile app (iOS/Android simulator required)
	@echo "Starting Flutter mobile app..."
	@echo "Note: Requires iOS Simulator or Android Emulator running"
	cd apps/mobile && flutter run

mobile-macos: ## Run Flutter app on macOS desktop
	@echo "Starting Flutter app on macOS..."
	cd apps/mobile && flutter run -d macos

mobile-web: ## Run Flutter app in Chrome browser
	@echo "Starting Flutter app in browser..."
	cd apps/mobile && flutter run -d chrome

mobile-setup: ## Enable macOS and Web platforms for Flutter
	@echo "Enabling macOS and Web platforms..."
	cd apps/mobile && flutter create --platforms=macos,web .
	@echo "✅ Platforms enabled! Now you can use:"
	@echo "  make mobile-macos  # Run on macOS"
	@echo "  make mobile-web    # Run in Chrome"

dev: docker-up ## Start all services in development mode
	@echo "Starting all services..."
	@echo "Database: http://localhost:5432"
	@echo "Backend API: http://localhost:8080"
	@echo "Web App: http://localhost:3000"
	@echo ""
	@echo "Run the following commands in separate terminals:"
	@echo "  make db-migrate-up # Apply DB migrations"
	@echo "  make server  # Start Go backend"
	@echo "  make web     # Start Next.js web"
	@echo "  make mobile  # Start Flutter app"

test: ## Run all tests (Go + Web + Mobile)
	@echo "Running Go tests..."
	cd apps/server && go test ./...
	@echo ""
	@echo "Running Next.js tests with Vitest..."
	cd apps/web && pnpm test
	@echo ""
	@echo "Running Flutter tests..."
	cd apps/mobile && flutter test

test-go: ## Run Go backend tests only
	@echo "Running Go tests..."
	cd apps/server && go test ./... -v

test-web: ## Run Next.js web tests only
	@echo "Running Next.js tests with Vitest..."
	cd apps/web && pnpm test

test-mobile: ## Run Flutter mobile tests only
	@echo "Running Flutter tests..."
	cd apps/mobile && flutter test

test-mobile-watch: ## Run Flutter tests in watch mode
	@echo "Running Flutter tests in watch mode..."
	cd apps/mobile && flutter test --watch

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf apps/web/.next
	rm -rf apps/web/node_modules
	rm -rf apps/mobile/build
	rm -rf node_modules
	@echo "Clean complete!"

reset: docker-down clean docker-up ## Reset everything (stop containers, clean, restart)
	@echo "Reset complete!"
