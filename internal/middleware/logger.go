package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
	"go.uber.org/zap"
)

// EchoLogger retorna un middleware de Echo que registra las peticiones
// En producción usa formato JSON, en otros ambientes usa formato legible
func Logger() echo.MiddlewareFunc {
	logger := logger.New()

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogMethod:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		HandleError:  true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			fields := []zap.Field{
				zap.String("method", values.Method),
				zap.String("uri", values.URI),
				zap.Int("status", values.Status),
				zap.Duration("latency", values.Latency),
				zap.String("remote_ip", values.RemoteIP),
				zap.String("user_agent", values.UserAgent),
			}

			if values.Error != nil {
				fields = append(fields, zap.Error(values.Error))
			}

			logger.Info("HTTP Request", fields...)
			return nil
		},
	})

	// Para desarrollo, usar el logger por defecto de Echo con formato legible
	// return middleware.Logger()
}
