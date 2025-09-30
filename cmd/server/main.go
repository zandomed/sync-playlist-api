package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zandomed/sync-playlist-api/internal/adapters/container"
	"github.com/zandomed/sync-playlist-api/internal/adapters/http/routes"
	"github.com/zandomed/sync-playlist-api/internal/config"
	SPMiddleware "github.com/zandomed/sync-playlist-api/internal/middleware"

	"github.com/zandomed/sync-playlist-api/pkg/database"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

func main() {
	cfg := config.Get()

	log := logger.New()
	defer log.Sync()

	if cfg.Server.Environment == config.Development {
		log.Sugar().Info("Running in development mode", cfg)
	}

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Sugar().Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	container := container.NewContainer(db, cfg, log)

	e := echo.New()
	e.HideBanner = true
	e.Validator = &SPMiddleware.CustomValidator{Validator: validator.New()}
	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	routes.SetupRoutes(e, container)

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      e,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Sugar().Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Sugar().Infof("Server started on port http://%s:%s", cfg.Server.Host, cfg.Server.Port)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Sugar().Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Sugar().Fatalf("Server forced to shutdown: %v", err)
	}

	log.Sugar().Info("Server exited")
}
