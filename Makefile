.PHONY: help build run test clean docker-up docker-down migrate dev commit-setup commit-validate commit-help setup

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
	@${DOCKER_COMPOSE} up -d
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

setup: ## Complete project setup (dependencies, tools, git hooks)
	@echo "${BLUE}🚀 Setting up ${APP_NAME} project...${NC}"
	@echo ""
	@$(MAKE) deps
	@$(MAKE) install-tools
	@$(MAKE) commit-setup
	@echo ""
	@echo "${GREEN}✅ Project setup completed!${NC}"
	@echo "${BLUE}You can now run:${NC}"
	@echo "  ${YELLOW}make dev${NC}       - Start development server"
	@echo "  ${YELLOW}make test${NC}      - Run tests"
	@echo "  ${YELLOW}make docker-up${NC} - Start services"

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
	#	github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
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

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "${RED}Error: NAME parameter is required${NC}"; \
		echo "${YELLOW}Usage: make migrate-create NAME=migration_name${NC}"; \
		echo "${YELLOW}Example: make migrate-create NAME=add_user_table${NC}"; \
		exit 1; \
	fi
	@echo "${YELLOW}Creating migration: $(NAME)${NC}"
	@go run scripts/generate_migration.go $(NAME)

commit-setup: ## Setup conventional commit hooks
	@echo "${YELLOW}Setting up conventional commit hooks...${NC}"
	@echo "${YELLOW}Installing commitlint dependencies...${NC}"
	@npm install
	@git config core.hooksPath .githooks
	@echo "${GREEN}✅ Git hooks configured${NC}"
	@echo "${BLUE}Hooks will now validate commit messages using commitlint${NC}"

commit-validate: ## Validate the last commit message
	@echo "${YELLOW}Validating last commit message...${NC}"
	@npm run commitlint-last

commit-help: ## Show conventional commit format help
	@echo "${BLUE}Conventional Commit Format:${NC}"
	@echo ""
	@echo "${YELLOW}Format:${NC} <type>[optional scope]: <description>"
	@echo ""
	@echo "${YELLOW}Types:${NC}"
	@echo "  ${GREEN}feat${NC}     - A new feature"
	@echo "  ${GREEN}fix${NC}      - A bug fix"
	@echo "  ${GREEN}docs${NC}     - Documentation only changes"
	@echo "  ${GREEN}style${NC}    - Changes that do not affect meaning (white-space, formatting, etc)"
	@echo "  ${GREEN}refactor${NC} - A code change that neither fixes a bug nor adds a feature"
	@echo "  ${GREEN}perf${NC}     - A code change that improves performance"
	@echo "  ${GREEN}test${NC}     - Adding missing tests or correcting existing tests"
	@echo "  ${GREEN}build${NC}    - Changes that affect the build system or external dependencies"
	@echo "  ${GREEN}ci${NC}       - Changes to CI configuration files and scripts"
	@echo "  ${GREEN}chore${NC}    - Other changes that don't modify src or test files"
	@echo "  ${GREEN}revert${NC}   - Reverts a previous commit"
	@echo ""
	@echo "${YELLOW}Examples:${NC}"
	@echo "  ${GREEN}feat(auth): add user authentication${NC}"
	@echo "  ${GREEN}fix(api): resolve validation error${NC}"
	@echo "  ${GREEN}docs: update README with setup instructions${NC}"

# Default values
.DEFAULT_GOAL := help