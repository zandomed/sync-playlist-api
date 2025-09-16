# Build stage
FROM golang:1.25.1-alpine3.22 AS builder

ARG PORT

# Instalar dependencias de build
RUN apk add --no-cache git ca-certificates tzdata

# Crear directorio de trabajo
WORKDIR /app

# Copiar módulos Go
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Runtime stage
FROM alpine:3.22.1

# Instalar certificados y zona horaria
RUN apk --no-cache add ca-certificates tzdata

# Crear usuario no-root
RUN adduser -D -s /bin/sh appuser

# Crear directorio de trabajo
WORKDIR /root/

# Copiar el binario desde el builder
COPY --from=builder /app/main .

# Cambiar propietario
RUN chown appuser:appuser main

# Cambiar a usuario no-root
USER appuser

# Variables de entorno por defecto
ENV PORT=8080
ENV HOST=0.0.0.0
ENV LOG_LEVEL=INFO

# Puerto de la aplicación
EXPOSE 8080

# Comando por defecto
CMD ["./main"]