package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Claims estructura personalizada para JWT
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// JWT middleware personalizado
func JWT(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Obtener token del header Authorization
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				// Intentar obtener del query parameter
				auth = c.QueryParam("token")
				if auth == "" {
					return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization token")
				}
			} else {
				// Remover "Bearer " del header
				if strings.HasPrefix(auth, "Bearer ") {
					auth = strings.TrimPrefix(auth, "Bearer ")
				} else {
					return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
				}
			}

			// Parsear y validar token
			token, err := jwt.ParseWithClaims(auth, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				// Verificar método de firma
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token signing method")
				}
				return []byte(secret), nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			// Verificar que el token es válido
			if !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			// Guardar token en el contexto
			c.Set("user", token)

			return next(c)
		}
	}
}

// GetUserFromContext extrae el usuario del contexto JWT
func GetUserFromContext(c echo.Context) (*Claims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token format")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
	}

	return claims, nil
}

// RequireAuth middleware que requiere autenticación
func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := GetUserFromContext(c)
		if err != nil {
			return err
		}
		return next(c)
	}
}
