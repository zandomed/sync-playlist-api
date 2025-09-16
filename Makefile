.PHONY: help build run test clean docker-up docker-down migrate dev

# Variables
APP_NAME=sync-playlist
BINARY_NAME=main
DOCKER_COMPOSE=docker-compose

# Colores para output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

help: ## Mostrar esta ayuda
	@echo "${BLUE}${APP_NAME} - Comandos disponibles:${NC}"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-15s${NC} %s\n", $$1, $$2}'

build: ## Compilar la aplicación
	@echo "${YELLOW}Compilando ${APP_NAME}...${NC}"
	@go build -o dist/${BINARY_NAME} cmd/server/main.go
	@echo "${GREEN}✅ Compilación completada${NC}"

run: build ## Compilar y ejecutar la aplicación
	@echo "${YELLOW}Ejecutando ${APP_NAME}...${NC}"
	@./dist/${BINARY_NAME}

dev: ## Ejecutar en modo desarrollo con live reload (requiere air)
	@echo "${YELLOW}Ejecutando en modo desarrollo...${NC}"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "${RED}Air no está instalado. Instalalo con: go install github.com/air-verse/air@latest${NC}"; \
		echo "${YELLOW}Ejecutando sin live reload...${NC}"; \
		go run cmd/server/main.go; \
	fi

test: ## Ejecutar tests
	@echo "${YELLOW}Ejecutando tests...${NC}"
	@go test -v ./...
	@echo "${GREEN}✅ Tests completados${NC}"

test-coverage: ## Ejecutar tests con cobertura
	@echo "${YELLOW}Ejecutando tests con cobertura...${NC}"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}✅ Reporte de cobertura generado en coverage.html${NC}"

clean: ## Limpiar archivos generados
	@echo "${YELLOW}Limpiando archivos...${NC}"
	@rm -f ${BINARY_NAME}
	@rm -f coverage.out coverage.html
	@echo "${GREEN}✅ Limpieza completada${NC}"

docker-up: ## Levantar servicios de desarrollo (PostgreSQL y Redis)
	@echo "${YELLOW}Levantando servicios de desarrollo...${NC}"
	@${DOCKER_COMPOSE} up -d postgres redis
	@echo "${GREEN}✅ Servicios levantados${NC}"

docker-full: ## Levantar todos los servicios incluyendo la app
	@echo "${YELLOW}Levantando todos los servicios...${NC}"
	@${DOCKER_COMPOSE} --profile full up -d
	@echo "${GREEN}✅ Todos los servicios levantados${NC}"

docker-down: ## Detener servicios de desarrollo
	@echo "${YELLOW}Deteniendo servicios...${NC}"
	@${DOCKER_COMPOSE} down
	@echo "${GREEN}✅ Servicios detenidos${NC}"

docker-logs: ## Ver logs de los servicios
	@${DOCKER_COMPOSE} logs -f

docker-clean: ## Limpiar volúmenes de Docker
	@echo "${YELLOW}Limpiando volúmenes de Docker...${NC}"
	@${DOCKER_COMPOSE} down -v
	@echo "${GREEN}✅ Volúmenes limpiados${NC}"

deps: ## Instalar dependencias
	@echo "${YELLOW}Instalando dependencias...${NC}"
	@go mod tidy
	@go mod download
	@echo "${GREEN}✅ Dependencias instaladas${NC}"

lint: ## Ejecutar linter (requiere golangci-lint)
	@echo "${YELLOW}Ejecutando linter...${NC}"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "${RED}golangci-lint no está instalado${NC}"; \
		echo "${YELLOW}Instálalo desde: https://golangci-lint.run/usage/install/${NC}"; \
	fi

format: ## Formatear código
	@echo "${YELLOW}Formateando código...${NC}"
	@go fmt ./...
	@echo "${GREEN}✅ Código formateado${NC}"

install-tools: ## Instalar herramientas de desarrollo
	@echo "${YELLOW}Instalando herramientas de desarrollo...${NC}"
	@go install github.com/air-verse/air@latest
	@echo "${GREEN}✅ Air instalado para live reload${NC}"

migrate-up: ## Ejecutar migraciones de base de datos
	@echo "${YELLOW}Ejecutando migraciones...${NC}"
	@go run cmd/migrate/main.go up
	@echo "${GREEN}✅ Migraciones completadas${NC}"

migrate-down: ## Rollback última migración
	@echo "${YELLOW}Haciendo rollback de migración...${NC}"
	@go run cmd/migrate/main.go down
	@echo "${GREEN}✅ Rollback completado${NC}"

migrate-status: ## Ver estado de migraciones
	@echo "${YELLOW}Estado de migraciones:${NC}"
	@go run cmd/migrate/main.go status

# Valores por defecto
.DEFAULT_GOAL := help