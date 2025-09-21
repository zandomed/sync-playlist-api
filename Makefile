.PHONY: help build run test clean docker-up docker-down migrate dev

# Variables
APP_NAME=sync-playlist
BINARY_NAME=main
DOCKER_COMPOSE=docker-compose

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

help: ## Show this help
	@echo "${BLUE}${APP_NAME} - Available commands:${NC}"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-15s${NC} %s\n", $$1, $$2}'

build: ## Build the application
	@echo "${YELLOW}Building ${APP_NAME}...${NC}"
	@go build -o dist/${BINARY_NAME} cmd/server/main.go
	@echo "${GREEN}✅ Build completed${NC}"

run: build ## Build and run the application
	@echo "${YELLOW}Running ${APP_NAME}...${NC}"
	@./dist/${BINARY_NAME}

dev: ## Run in development mode with live reload (requires air)
	@echo "${YELLOW}Running in development mode...${NC}"
	@if command -v air > /dev/null; then \
		air -c ./.air.toml; \
	else \
		echo "${RED}Air is not installed. Install it with: go install github.com/air-verse/air@latest${NC}"; \
		echo "${YELLOW}Running without live reload...${NC}"; \
		go run cmd/server/main.go; \
	fi

test: ## Run tests
	@echo "${YELLOW}Running tests...${NC}"
	@go test -v ./...
	@echo "${GREEN}✅ Tests completed${NC}"

test-coverage: ## Run tests with coverage
	@echo "${YELLOW}Running tests with coverage...${NC}"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}✅ Coverage report generated in coverage.html${NC}"

clean: ## Clean generated files
	@echo "${YELLOW}Cleaning files...${NC}"
	@rm -f coverage.out coverage.html
	@rm -rf dist/ tmp/
	@echo "${GREEN}✅ Cleanup completed${NC}"

docker-up: ## Start development services (PostgreSQL and Redis)
	@echo "${YELLOW}Starting development services...${NC}"
	@${DOCKER_COMPOSE} up -d postgres redis
	@echo "${GREEN}✅ Services started${NC}"

docker-full: ## Start all services including the app
	@echo "${YELLOW}Starting all services...${NC}"
	@${DOCKER_COMPOSE} --profile full up -d
	@echo "${GREEN}✅ All services started${NC}"

docker-down: ## Stop development services
	@echo "${YELLOW}Stopping services...${NC}"
	@${DOCKER_COMPOSE} down
	@echo "${GREEN}✅ Services stopped${NC}"

docker-logs: ## View service logs
	@${DOCKER_COMPOSE} logs -f

docker-clean: ## Clean Docker volumes
	@echo "${YELLOW}Cleaning Docker volumes...${NC}"
	@${DOCKER_COMPOSE} down -v
	@echo "${GREEN}✅ Volumes cleaned${NC}"

deps: ## Install dependencies
	@echo "${YELLOW}Installing dependencies...${NC}"
	@go mod tidy
	@go mod download
	@echo "${GREEN}✅ Dependencies installed${NC}"

lint: ## Run linter (requires golangci-lint)
	@echo "${YELLOW}Running linter...${NC}"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "${RED}golangci-lint is not installed${NC}"; \
		echo "${YELLOW}Install it from: https://golangci-lint.run/usage/install/${NC}"; \
	fi

format: ## Format code
	@echo "${YELLOW}Formatting code...${NC}"
	@if command -v golangci-lint > /dev/null; then \
		echo "${YELLOW}Using golangci-lint fmt...${NC}"; \
		golangci-lint fmt; \
	else \
		echo "${RED}golangci-lint is not installed${NC}"; \
		echo "${YELLOW}Install it from: https://golangci-lint.run/usage/install/${NC}"; \
		echo "${YELLOW}Running go fmt instead...${NC}"; \
		go fmt ./...; \
	fi
	@echo "${GREEN}✅ Code formatted${NC}"

doctor: ## Check installed tools
	@echo "${BLUE}🔍 Checking installed tools:${NC}"
	@echo ""
	@echo "${BLUE}Go:${NC}"
	@if command -v go > /dev/null; then \
		echo "  ${GREEN}✅ $$(go version)${NC}"; \
	else \
		echo "  ${RED}❌ Go not found${NC}"; \
	fi
	@echo "${BLUE}Docker:${NC}"
	@if command -v docker > /dev/null; then \
		echo "  ${GREEN}✅ $$(docker --version)${NC}"; \
	else \
		echo "  ${RED}❌ Docker not found${NC}"; \
	fi
	@echo "${BLUE}Docker Compose:${NC}"
	@if command -v docker-compose > /dev/null; then \
		echo "  ${GREEN}✅ $$(docker-compose --version)${NC}"; \
	else \
		echo "  ${RED}❌ Docker Compose not found${NC}"; \
	fi
	@echo "${BLUE}Air (live reload):${NC}"
	@if command -v air > /dev/null; then \
		echo "  ${GREEN}✅ Air found${NC}"; \
	else \
		echo "  ${YELLOW}❌ Air not found (run: go install github.com/air-verse/air@latest)${NC}"; \
	fi
	@echo "${BLUE}Golangci-lint:${NC}"
	@if command -v golangci-lint > /dev/null; then \
		echo "  ${GREEN}✅ Golangci-lint found${NC}"; \
	else \
		echo "  ${YELLOW}❌ Golangci-lint not found${NC}"; \
		echo "    ${BLUE}Download from: https://golangci-lint.run/docs/welcome/install/${NC}"; \
	fi

install-tools: ## Install development tools
	@echo "${YELLOW}Installing development tools...${NC}"
	@go install github.com/air-verse/air@latest
	@echo "${GREEN}✅ Air installed for live reload${NC}"

migrate-up: ## Run database migrations
	@echo "${YELLOW}Running migrations...${NC}"
	@go run cmd/migrate/main.go up
	@echo "${GREEN}✅ Migrations completed${NC}"

migrate-down: ## Rollback last migration
	@echo "${YELLOW}Rolling back migration...${NC}"
	@go run cmd/migrate/main.go down
	@echo "${GREEN}✅ Rollback completed${NC}"

migrate-status: ## View migration status
	@echo "${YELLOW}Migration status:${NC}"
	@go run cmd/migrate/main.go status

# Default values
.DEFAULT_GOAL := help