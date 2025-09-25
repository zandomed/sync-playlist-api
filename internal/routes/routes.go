package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zandomed/sync-playlist-api/internal/handlers"
)

func Setup(e *echo.Echo, h *handlers.Handlers) {

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	api := e.Group("/v1")

	auth := api.Group("/auth")
	{
		auth.POST("/login", h.Auth.LoginWithPass)
		auth.POST("/register", h.Auth.Register)
	}
}
