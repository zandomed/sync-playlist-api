package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zandomed/sync-playlist-api/internal/config"
)

func CORS() echo.MiddlewareFunc {
	cfg := config.Get()

	allowOrigins := []string{"http://localhost:3000", "http://localhost:5173"}
	if cfg.Server.Environment == config.Production {
		// Add your production frontend URLs here
		allowOrigins = []string{"https://yourdomain.com"}
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
		ExposeHeaders:    []string{echo.HeaderContentLength, echo.HeaderContentType},
		MaxAge:           86400, // 24 hours
	})
}