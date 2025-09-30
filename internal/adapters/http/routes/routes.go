package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/zandomed/sync-playlist-api/internal/adapters/container"
	"github.com/zandomed/sync-playlist-api/internal/middleware"
)

func SetupRoutes(e *echo.Echo, container *container.Container) {
	e.Use(middleware.CORS())

	e.GET("/health", container.HealthHandler.GetStatus)
	e.GET("/", container.HealthHandler.GetStatus)
	api := e.Group("/v1")

	api.Use(middleware.Logger())

	auth := api.Group("/auth")
	{
		auth.POST("/register", container.AuthHandler.Register)
		auth.POST("/login", container.AuthHandler.Login)
		auth.GET("/google", container.AuthHandler.GoogleAuth)
		auth.GET("/google/callback", container.AuthHandler.GoogleCallback)
		auth.GET("/spotify", container.AuthHandler.SpotifyAuth)
		auth.GET("/spotify/callback", container.AuthHandler.SpotifyCallback)
	}
}
