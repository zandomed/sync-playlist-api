package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/handlers"
	custMiddleware "github.com/zandomed/sync-playlist-api/internal/middleware"
)

func Setup(e *echo.Echo, h *handlers.Handlers, cfg *config.Config) {
	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	// API v1
	api := e.Group("/api/v1")

	// Auth routes (sin autenticación requerida)
	auth := api.Group("/auth")
	auth.GET("/spotify", h.Auth.SpotifyAuth)
	auth.GET("/spotify/callback", h.Auth.SpotifyCallback)
	auth.GET("/apple", h.Auth.AppleAuth)
	auth.GET("/apple/callback", h.Auth.AppleCallback)
	auth.POST("/refresh", h.Auth.RefreshToken)

	// Protected routes (requieren autenticación)
	protected := api.Group("")
	protected.Use(custMiddleware.JWT(cfg.JWT.Secret))

	// User routes
	users := protected.Group("/users")
	users.GET("/me", h.User.GetProfile)
	users.PUT("/me", h.User.UpdateProfile)

	// Playlists routes
	playlists := protected.Group("/playlists")
	playlists.GET("", h.Playlist.GetPlaylists)
	playlists.GET("/:id", h.Playlist.GetPlaylist)

	// Migration routes
	migrations := protected.Group("/migrations")
	migrations.POST("", h.Migration.StartMigration)
	migrations.GET("", h.Migration.GetUserMigrations)
	migrations.GET("/:id", h.Migration.GetMigrationStatus)
	migrations.GET("/:id/progress", h.Migration.GetMigrationProgress)
	migrations.DELETE("/:id", h.Migration.CancelMigration)

	// WebSocket para progreso en tiempo real
	e.GET("/ws/migration/:id", h.Migration.WebSocketHandler)
}
