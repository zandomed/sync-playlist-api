package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/zandomed/sync-playlist-api/internal/infra/container"
	"github.com/zandomed/sync-playlist-api/internal/middleware"
)

func SetupRoutes(e *echo.Echo, container *container.Container) {
	e.Use(middleware.CORS())

	e.GET("/health", container.HealthHandler.GetStatus)
	e.GET("/", container.HealthHandler.GetStatus)
	api := e.Group("/v1")

	api.Use(middleware.Logger())

	oauth := api.Group("/oauth")
	{
		oauth.GET("/google", container.AuthHandler.GoogleAuth)
		oauth.GET("/google/callback", container.AuthHandler.GoogleCallback)
		oauth.GET("/spotify", container.AuthHandler.SpotifyAuth)
		oauth.GET("/spotify/callback", container.AuthHandler.SpotifyCallback)
		oauth.POST("/verify", container.AuthHandler.VerifyToken)
	}

	auth := api.Group("/auth")
	{
		auth.POST("/register", container.AuthHandler.Register)
		auth.POST("/login", container.AuthHandler.Login)
	}
}
