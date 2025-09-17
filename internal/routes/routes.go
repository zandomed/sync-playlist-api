package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo) {
	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})
}
