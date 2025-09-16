package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/handlers"
	SPMiddleware "github.com/zandomed/sync-playlist-api/internal/middleware"
	"github.com/zandomed/sync-playlist-api/internal/repository"
	"github.com/zandomed/sync-playlist-api/internal/routes"
	"github.com/zandomed/sync-playlist-api/internal/services"
	"github.com/zandomed/sync-playlist-api/pkg/database"
	SPLogger "github.com/zandomed/sync-playlist-api/pkg/logger"
)

func main() {
	// Obtener configuración singleton
	cfg := config.Get()

	// Inicializar logger
	logger := SPLogger.New()

	// Conectar a la base de datos
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		logger.Sugar().Fatalf("Error connecting to database: %v", err)
	}

	// Inicializar Echo
	e := echo.New()

	// Configurar Echo
	e.HideBanner = true
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(SPMiddleware.Logger())

	// Inicializar dependencias
	repos := repository.New(db)
	services := services.New(repos, cfg)
	handlers := handlers.New(services, cfg, logger)

	// Configurar rutas
	routes.Setup(e, handlers, cfg)

	// Servidor con graceful shutdown
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      e,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Iniciar servidor en goroutine
	go func() {
		logger.Info(fmt.Sprintf("🚀 Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port))
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			logger.Sugar().Fatalf("Error starting server: %v", err)
		}
	}()

	// Esperar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger.Info("🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Sugar().Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("✅ Server stopped gracefully")
}
