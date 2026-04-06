.PHONY: help build up down restart logs clean migrate-up migrate-down migrate-reset db-shell redis-shell test lint

# Variables
DOCKER_COMPOSE = docker compose
DB_USER = dnsproxy
DB_PASSWORD = dev_password_123
DB_NAME = dnsproxy
DB_HOST = localhost
DB_PORT = 5432
REDIS_HOST = localhost
REDIS_PORT = 6379

# Colors for output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

## Help
help: ## Show this help message
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${WHITE}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

## Development
build: ## Build all Docker images
	@echo "${GREEN}Building Docker images...${RESET}"
	$(DOCKER_COMPOSE) build

up: ## Start all services (detached)
	@echo "${GREEN}Starting all services...${RESET}"
	$(DOCKER_COMPOSE) up -d
	@echo "${GREEN}Services started successfully!${RESET}"
	@echo "${YELLOW}PostgreSQL:${RESET} $(DB_HOST):$(DB_PORT)"
	@echo "${YELLOW}Redis:${RESET} $(REDIS_HOST):$(REDIS_PORT)"

up-infra: ## Start only infrastructure (PostgreSQL + Redis)
	@echo "${GREEN}Starting infrastructure services...${RESET}"
	$(DOCKER_COMPOSE) up -d postgres redis
	@echo "${GREEN}Infrastructure started successfully!${RESET}"

down: ## Stop all services
	@echo "${YELLOW}Stopping all services...${RESET}"
	$(DOCKER_COMPOSE) down

down-volumes: ## Stop all services and remove volumes
	@echo "${YELLOW}Stopping all services and removing volumes...${RESET}"
	$(DOCKER_COMPOSE) down -v

restart: down up ## Restart all services

restart-infra: ## Restart only infrastructure
	@echo "${YELLOW}Restarting infrastructure...${RESET}"
	$(DOCKER_COMPOSE) restart postgres redis

logs: ## Show logs from all services
	$(DOCKER_COMPOSE) logs -f

logs-dns: ## Show logs from DNS server
	$(DOCKER_COMPOSE) logs -f dns-server

logs-proxy: ## Show logs from Proxy server
	$(DOCKER_COMPOSE) logs -f proxy-server

logs-api: ## Show logs from API server
	$(DOCKER_COMPOSE) logs -f api-server

logs-bot: ## Show logs from Telegram bot
	$(DOCKER_COMPOSE) logs -f telegram-bot

logs-db: ## Show logs from PostgreSQL
	$(DOCKER_COMPOSE) logs -f postgres

logs-redis: ## Show logs from Redis
	$(DOCKER_COMPOSE) logs -f redis

ps: ## Show running containers
	$(DOCKER_COMPOSE) ps

## Database
migrate-up: ## Apply all database migrations
	@echo "${GREEN}Applying migrations...${RESET}"
	@./scripts/migrate.sh
	@echo "${GREEN}Migrations applied successfully!${RESET}"

migrate-down: ## Rollback last migration
	@echo "${YELLOW}Rolling back last migration...${RESET}"
	@for file in $$(ls -r migrations/*.down.sql | head -1); do \
		echo "Applying: $$file"; \
		$(DOCKER_COMPOSE) exec -T postgres psql -U $(DB_USER) -d $(DB_NAME) < "$$file"; \
	done
	@echo "${GREEN}Migration rolled back successfully!${RESET}"

migrate-reset: ## Reset database (drop all tables and reapply migrations)
	@echo "${YELLOW}Resetting database...${RESET}"
	@for file in $$(ls -r migrations/*.down.sql); do \
		echo "Applying: $$file"; \
		$(DOCKER_COMPOSE) exec -T postgres psql -U $(DB_USER) -d $(DB_NAME) < "$$file" 2>/dev/null || true; \
	done
	@echo "${GREEN}Database reset. Applying migrations...${RESET}"
	@$(MAKE) migrate-up

db-shell: ## Open PostgreSQL shell
	@echo "${GREEN}Opening PostgreSQL shell...${RESET}"
	$(DOCKER_COMPOSE) exec postgres psql -U $(DB_USER) -d $(DB_NAME)

db-tables: ## List all database tables
	@echo "${GREEN}Database tables:${RESET}"
	@$(DOCKER_COMPOSE) exec postgres psql -U $(DB_USER) -d $(DB_NAME) -c "\dt"

db-clean: ## Drop all tables
	@echo "${YELLOW}Dropping all tables...${RESET}"
	@for file in $$(ls -r migrations/*.down.sql); do \
		$(DOCKER_COMPOSE) exec -T postgres psql -U $(DB_USER) -d $(DB_NAME) < "$$file" 2>/dev/null || true; \
	done
	@echo "${GREEN}All tables dropped!${RESET}"

db-seed: ## Seed database with test data
	@echo "${GREEN}Seeding database...${RESET}"
	@if [ -f scripts/seed.sql ]; then \
		$(DOCKER_COMPOSE) exec -T postgres psql -U $(DB_USER) -d $(DB_NAME) < scripts/seed.sql; \
		echo "${GREEN}Database seeded successfully!${RESET}"; \
	else \
		echo "${YELLOW}No seed file found at scripts/seed.sql${RESET}"; \
	fi

redis-shell: ## Open Redis CLI
	@echo "${GREEN}Opening Redis CLI...${RESET}"
	$(DOCKER_COMPOSE) exec redis redis-cli

redis-flush: ## Flush all Redis data
	@echo "${YELLOW}Flushing Redis...${RESET}"
	@$(DOCKER_COMPOSE) exec redis redis-cli FLUSHALL
	@echo "${GREEN}Redis flushed!${RESET}"

## Application
run-dns: ## Run DNS server locally
	@echo "${GREEN}Starting DNS server...${RESET}"
	go run cmd/dns-server/main.go

run-proxy: ## Run Proxy server locally
	@echo "${GREEN}Starting Proxy server...${RESET}"
	go run cmd/proxy-server/main.go

run-api: ## Run API server locally
	@echo "${GREEN}Starting API server...${RESET}"
	go run cmd/api-server/main.go

run-bot: ## Run Telegram bot locally
	@echo "${GREEN}Starting Telegram bot...${RESET}"
	go run cmd/telegram-bot/main.go

## Build
build-dns: ## Build DNS server binary
	@echo "${GREEN}Building DNS server...${RESET}"
	go build -o bin/dns-server cmd/dns-server/main.go

build-proxy: ## Build Proxy server binary
	@echo "${GREEN}Building Proxy server...${RESET}"
	go build -o bin/proxy-server cmd/proxy-server/main.go

build-api: ## Build API server binary
	@echo "${GREEN}Building API server...${RESET}"
	go build -o bin/api-server cmd/api-server/main.go

build-bot: ## Build Telegram bot binary
	@echo "${GREEN}Building Telegram bot...${RESET}"
	go build -o bin/telegram-bot cmd/telegram-bot/main.go

build-all: build-dns build-proxy build-api build-bot ## Build all binaries

## Testing
test: ## Run all tests
	@echo "${GREEN}Running tests...${RESET}"
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@echo "${GREEN}Generating coverage report...${RESET}"
	go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}Coverage report generated: coverage.html${RESET}"

## Code Quality
lint: ## Run linter
	@echo "${GREEN}Running linter...${RESET}"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "${YELLOW}golangci-lint not installed. Install it with:${RESET}"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt: ## Format code
	@echo "${GREEN}Formatting code...${RESET}"
	go fmt ./...

vet: ## Run go vet
	@echo "${GREEN}Running go vet...${RESET}"
	go vet ./...

tidy: ## Tidy go modules
	@echo "${GREEN}Tidying go modules...${RESET}"
	go mod tidy

## Cleanup
clean: down-volumes ## Clean up everything (containers, volumes, binaries)
	@echo "${YELLOW}Cleaning up...${RESET}"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "${GREEN}Cleanup complete!${RESET}"

clean-cache: ## Clean Go build cache
	@echo "${YELLOW}Cleaning Go cache...${RESET}"
	go clean -cache -testcache -modcache
	@echo "${GREEN}Cache cleaned!${RESET}"

## Setup
setup: ## Initial setup (install dependencies, start infra, run migrations)
	@echo "${GREEN}Setting up project...${RESET}"
	@$(MAKE) tidy
	@$(MAKE) up-infra
	@sleep 5
	@$(MAKE) migrate-up
	@echo "${GREEN}Setup complete!${RESET}"

## Info
info: ## Show project information
	@echo "${GREEN}Project Information:${RESET}"
	@echo "${YELLOW}PostgreSQL:${RESET}"
	@echo "  Host: $(DB_HOST)"
	@echo "  Port: $(DB_PORT)"
	@echo "  User: $(DB_USER)"
	@echo "  Database: $(DB_NAME)"
	@echo ""
	@echo "${YELLOW}Redis:${RESET}"
	@echo "  Host: $(REDIS_HOST)"
	@echo "  Port: $(REDIS_PORT)"
	@echo ""
	@echo "${YELLOW}Services:${RESET}"
	@echo "  DNS Server: port 53 (UDP/TCP)"
	@echo "  Proxy HTTP: port 8080"
	@echo "  Proxy HTTPS: port 8443"
	@echo "  API Server: port 8000"
