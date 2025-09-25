package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/services"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

// import (
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/google/uuid"
// 	"github.com/labstack/echo/v4"

// 	"github.com/zandomed/sync-playlist-api/internal/config"
// 	custMiddleware "github.com/zandomed/sync-playlist-api/internal/middleware"
// 	"github.com/zandomed/sync-playlist-api/internal/services"
// 	"github.com/zandomed/sync-playlist-api/pkg/logger"
// )

// // Handlers contiene todos los handlers HTTP
type Handlers struct {
	Auth *AuthHandler
	// User      *UserHandler
	// Playlist  *PlaylistHandler
	// Migration *MigrationHandler
}

// New crea una nueva instancia de handlers
func New(services *services.Services, cfg *config.Config, logger *logger.Logger) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(services.Auth, cfg, logger),
		// User:      NewUserHandler(services.User, logger),
		// Playlist:  NewPlaylistHandler(services.Playlist, logger),
		// Migration: NewMigrationHandler(services.Migration, logger),
	}
}

// Response structs comunes
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// // Utility functions
// func getUserFromContext(c echo.Context) (*custMiddleware.Claims, error) {
// 	return custMiddleware.GetUserFromContext(c)
// }

func sendError(c echo.Context, status int, err error, message string) error {
	response := ErrorResponse{
		Error:   err.Error(),
		Message: message,
	}
	return c.JSON(status, response)
}

func sendValidationError(c echo.Context, status int, errors error) error {
	response := ErrorResponse{
		Error:   "Validation error",
		Message: "Validation failed",
		Details: errors.Error(),
	}
	return c.JSON(status, response)
}

func sendSuccess(c echo.Context, status int, data interface{}, message string) error {
	response := SuccessResponse{
		Data:    data,
		Message: message,
	}
	return c.JSON(status, response)
}

// // AuthHandler maneja la autenticación
// type AuthHandler struct {
// 	authService services.AuthService
// 	userService services.UserService
// 	config      *config.Config
// 	logger      *logger.Logger
// }

// func NewAuthHandler(authService services.AuthService, userService services.UserService, cfg *config.Config, logger *logger.Logger) *AuthHandler {
// 	return &AuthHandler{
// 		authService: authService,
// 		userService: userService,
// 		config:      cfg,
// 		logger:      logger,
// 	}
// }

// // SpotifyAuth inicia el flujo OAuth de Spotify
// func (h *AuthHandler) SpotifyAuth(c echo.Context) error {
// 	// TODO: Implementar OAuth real de Spotify
// 	// Por ahora retornamos placeholder
// 	authURL := "https://accounts.spotify.com/authorize?client_id=" + h.config.Spotify.ClientID + "&response_type=code&redirect_uri=" + h.config.Spotify.RedirectURL + "&scope=playlist-read-private playlist-read-collaborative"

// 	return c.JSON(http.StatusOK, map[string]string{
// 		"auth_url": authURL,
// 	})
// }

// // SpotifyCallback maneja el callback de OAuth de Spotify
// func (h *AuthHandler) SpotifyCallback(c echo.Context) error {
// 	code := c.QueryParam("code")
// 	if code == "" {
// 		return sendError(c, http.StatusBadRequest, nil, "Authorization code required")
// 	}

// 	// TODO: Intercambiar código por token con Spotify
// 	// Por ahora simulamos un token

// 	// Crear o obtener usuario (simulado)
// 	user, err := h.userService.CreateUser("user@example.com", "Test User")
// 	if err != nil {
// 		return sendError(c, http.StatusInternalServerError, err, "Failed to create user")
// 	}

// 	// Guardar auth (simulado)
// 	expiresAt := time.Now().Add(1 * time.Hour)
// 	if err := h.authService.SaveAuth(user.ID, "spotify", "fake_access_token", "fake_refresh_token", expiresAt); err != nil {
// 		return sendError(c, http.StatusInternalServerError, err, "Failed to save auth")
// 	}

// 	// Crear JWT
// 	token, err := h.createJWT(user.ID, user.Email)
// 	if err != nil {
// 		return sendError(c, http.StatusInternalServerError, err, "Failed to create token")
// 	}

// 	return sendSuccess(c, http.StatusOK, map[string]interface{}{
// 		"token": token,
// 		"user":  user,
// 	}, "Authentication successful")
// }

// // AppleAuth inicia el flujo OAuth de Apple Music
// func (h *AuthHandler) AppleAuth(c echo.Context) error {
// 	// TODO: Implementar OAuth real de Apple Music
// 	authURL := "https://music.apple.com/auth?client_id=" + h.config.Apple.TeamID

// 	return c.JSON(http.StatusOK, map[string]string{
// 		"auth_url": authURL,
// 	})
// }

// // AppleCallback maneja el callback de OAuth de Apple Music
// func (h *AuthHandler) AppleCallback(c echo.Context) error {
// 	// TODO: Implementar callback real de Apple Music
// 	return sendError(c, http.StatusNotImplemented, nil, "Apple Music OAuth not implemented yet")
// }

// // RefreshToken renueva un token JWT
// func (h *AuthHandler) RefreshToken(c echo.Context) error {
// 	// TODO: Implementar refresh de token
// 	return sendError(c, http.StatusNotImplemented, nil, "Token refresh not implemented yet")
// }

// func (h *AuthHandler) createJWT(userID uuid.UUID, email string) (string, error) {
// 	claims := &custMiddleware.Claims{
// 		UserID: userID,
// 		Email:  email,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.config.JWT.ExpirationTime)),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString([]byte(h.config.JWT.Secret))
// }

// // UserHandler maneja operaciones de usuario
// type UserHandler struct {
// 	userService services.UserService
// 	logger      *logger.Logger
// }

// func NewUserHandler(userService services.UserService, logger *logger.Logger) *UserHandler {
// 	return &UserHandler{
// 		userService: userService,
// 		logger:      logger,
// 	}
// }

// // GetProfile obtiene el perfil del usuario actual
// func (h *UserHandler) GetProfile(c echo.Context) error {
// 	claims, err := getUserFromContext(c)
// 	if err != nil {
// 		return sendError(c, http.StatusUnauthorized, err, "Invalid token")
// 	}

// 	user, err := h.userService.GetUser(claims.UserID)
// 	if err != nil {
// 		return sendError(c, http.StatusNotFound, err, "User not found")
// 	}

// 	return sendSuccess(c, http.StatusOK, user, "")
// }

// // UpdateProfile actualiza el perfil del usuario
// func (h *UserHandler) UpdateProfile(c echo.Context) error {
// 	claims, err := getUserFromContext(c)
// 	if err != nil {
// 		return sendError(c, http.StatusUnauthorized, err, "Invalid token")
// 	}

// 	type UpdateProfileRequest struct {
// 		Name string `json:"name" validate:"required,min=2,max=100"`
// 	}

// 	req := new(UpdateProfileRequest)
// 	if err := c.Bind(req); err != nil {
// 		return sendError(c, http.StatusBadRequest, err, "Invalid request body")
// 	}

// 	if err := c.Validate(req); err != nil {
// 		return sendError(c, http.StatusBadRequest, err, "Validation failed")
// 	}

// 	user, err := h.userService.GetUser(claims.UserID)
// 	if err != nil {
// 		return sendError(c, http.StatusNotFound, err, "User not found")
// 	}

// 	user.Name = req.Name
// 	if err := h.userService.UpdateUser(user); err != nil {
// 		return sendError(c, http.StatusInternalServerError, err, "Failed to update user")
// 	}

// 	return sendSuccess(c, http.StatusOK, user, "Profile updated successfully")
// }

// // PlaylistHandler maneja operaciones de playlist
// type PlaylistHandler struct {
// 	playlistService services.PlaylistService
// 	logger          *logger.Logger
// }

// func NewPlaylistHandler(playlistService services.PlaylistService, logger *logger.Logger) *PlaylistHandler {
// 	return &PlaylistHandler{
// 		playlistService: playlistService,
// 		logger:          logger,
// 	}
// }

// // GetPlaylists obtiene las playlists del usuario
// func (h *PlaylistHandler) GetPlaylists(c echo.Context) error {
// 	claims, err := getUserFromContext(c)
// 	if err != nil {
// 		return sendError(c, http.StatusUnauthorized, err, "Invalid token")
// 	}

// 	// Parámetros de query
// 	service := c.QueryParam("service")
// 	pageStr := c.QueryParam("page")
// 	limitStr := c.QueryParam("limit")

// 	page := 1
// 	if pageStr != "" {
// 		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
// 			page = p
// 		}
// 	}

// 	limit := 20
// 	if limitStr != "" {
// 		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
// 			limit = l
// 		}
// 	}

// 	playlists, err := h.playlistService.GetUserPlaylists(claims.UserID, service, page, limit)
// 	if err != nil {
// 		return sendError(c, http.StatusInternalServerError, err, "Failed to get playlists")
// 	}

// 	return sendSuccess(c, http.StatusOK, map[string]interface{}{
// 		"playlists": playlists,
// 		"page":      page,
// 		"limit":     limit,
// 	}, "")
// }

// // GetPlaylist obtiene una playlist específica
// func (h *PlaylistHandler) GetPlaylist(c echo.Context) error {
// 	claims, err := getUserFromContext(c)
// 	if err != nil {
// 		return sendError(c, http.StatusUnauthorized, err, "Invalid token")
// 	}

// 	playlistIDStr := c.Param("id")
// 	playlistID, err := uuid.Parse(playlistIDStr)
// 	if err != nil {
// 		return sendError(c, http.StatusBadRequest, err, "Invalid playlist ID")
// 	}

// 	playlist, err := h.playlistService.GetPlaylist(playlistID)
// 	if err != nil {
// 		return sendError(c, http.StatusNotFound, err, "Playlist not found")
// 	}

// 	// Verificar que la playlist pertenece al usuario
// 	if playlist.UserID != claims.UserID {
// 		return sendError(c, http.StatusForbidden, nil, "Access denied")
// 	}

// 	// Incluir tracks si se solicita
// 	includeTracks := c.QueryParam("include_tracks") == "true"
// 	if includeTracks {
// 		tracks, err := h.playlistService.GetPlaylistTracks(playlistID)
// 		if err != nil {
// 			return sendError(c, http.StatusInternalServerError, err, "Failed to get tracks")
// 		}

// 		return sendSuccess(c, http.StatusOK, map[string]interface{}{
// 			"playlist": playlist,
// 			"tracks":   tracks,
// 		}, "")
// 	}

// 	return sendSuccess(c, http.StatusOK, playlist, "")
// }
