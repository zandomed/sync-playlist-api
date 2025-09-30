package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zandomed/sync-playlist-api/internal/adapters/http/dtos"
	"github.com/zandomed/sync-playlist-api/internal/adapters/http/mappers"
	"github.com/zandomed/sync-playlist-api/internal/middleware"
	"github.com/zandomed/sync-playlist-api/internal/usecases/auth"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

type AuthHandler struct {
	registerUseCase     *auth.RegisterUserUseCase
	loginUseCase        *auth.LoginUserUseCase
	googleLoginUseCase  *auth.GoogleLoginUseCase
	spotifyLoginUseCase *auth.SpotifyLoginUseCase
	linkSpotifyUseCase  *auth.LinkSpotifyAccountUseCase
	mapper              *mappers.AuthMapper
	logger              *logger.Logger
}

func NewAuthHandler(
	registerUseCase *auth.RegisterUserUseCase,
	loginUseCase *auth.LoginUserUseCase,
	googleLoginUseCase *auth.GoogleLoginUseCase,
	spotifyLoginUseCase *auth.SpotifyLoginUseCase,
	linkSpotifyUseCase *auth.LinkSpotifyAccountUseCase,
	mapper *mappers.AuthMapper,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		registerUseCase:     registerUseCase,
		loginUseCase:        loginUseCase,
		googleLoginUseCase:  googleLoginUseCase,
		spotifyLoginUseCase: spotifyLoginUseCase,
		linkSpotifyUseCase:  linkSpotifyUseCase,
		mapper:              mapper,
		logger:              logger,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var dto dtos.RegisterRequest
	if err := c.Bind(&dto); err != nil {
		h.logger.Sugar().Warnf("Invalid request body: %v", err)
		return SendError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := c.Validate(&dto); err != nil {
		h.logger.Sugar().Warnf("Validation failed: %v", err)
		return SendValidationError(c, err)
	}

	request := h.mapper.ToRegisterUserRequest(&dto)

	response, err := h.registerUseCase.Execute(c.Request().Context(), *request)
	if err != nil {
		h.logger.Sugar().Errorf("Registration failed: %v", err)
		return HandleUseCaseError(c, err)
	}

	h.logger.Sugar().Infof("User registered successfully: %s", response.UserID)
	return SendSuccess(c, http.StatusCreated, h.mapper.ToRegisterResponse(response))
}

func (h *AuthHandler) Login(c echo.Context) error {
	var dto dtos.LoginRequest
	if err := c.Bind(&dto); err != nil {
		h.logger.Sugar().Warnf("Invalid request body: %v", err)
		return SendError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := c.Validate(&dto); err != nil {
		h.logger.Sugar().Warnf("Validation failed: %v", err)
		return SendValidationError(c, err)
	}

	request := h.mapper.ToLoginUserRequest(&dto)

	response, err := h.loginUseCase.Execute(c.Request().Context(), *request)
	if err != nil {
		h.logger.Sugar().Warnf("Login failed for user %s: %v", dto.Email, err)
		return HandleUseCaseError(c, err)
	}

	h.logger.Sugar().Infof("User logged in successfully: %s", response.UserID)
	return SendSuccess(c, http.StatusOK, h.mapper.ToLoginResponse(response))
}

func (h *AuthHandler) GoogleAuth(c echo.Context) error {
	var dto dtos.GoogleAuthURLRequest
	if err := c.Bind(&dto); err != nil {
		dto.State = "random-state"
	}

	if dto.State == "" {
		dto.State = "random-state"
	}

	request := h.mapper.ToGoogleAuthURLRequest(&dto)

	response, err := h.googleLoginUseCase.GetAuthURL(c.Request().Context(), *request)
	if err != nil {
		h.logger.Sugar().Errorf("Failed to generate Google auth URL: %v", err)
		return HandleUseCaseError(c, err)
	}

	h.logger.Sugar().Infof("Generated Google auth URL")

	// Check if client accepts JSON
	acceptHeader := c.Request().Header.Get("Accept")
	if acceptHeader == string(ContentTypeApplicationJson) || c.Request().Header.Get("Content-Type") == string(ContentTypeApplicationJson) {
		return SendSuccess(c, http.StatusOK, h.mapper.ToGoogleAuthURLResponse(response))
	}
	return c.Redirect(http.StatusFound, response.URL)
}

func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		h.logger.Sugar().Warn("Missing authorization code in Google callback")
		return SendError(c, http.StatusBadRequest, "invalid_request", "Missing authorization code")
	}

	dto := dtos.GoogleCallbackRequest{
		Code:  code,
		State: state,
	}

	request := h.mapper.ToGoogleCallbackRequest(&dto)

	response, err := h.googleLoginUseCase.HandleCallback(c.Request().Context(), *request)
	if err != nil {
		h.logger.Sugar().Errorf("Google callback failed: %v", err)
		return HandleUseCaseError(c, err)
	}

	if response.IsNewUser {
		h.logger.Sugar().Infof("New user created via Google OAuth: %s", response.UserID)
	} else {
		h.logger.Sugar().Infof("User logged in via Google OAuth: %s", response.UserID)
	}

	return SendSuccess(c, http.StatusOK, h.mapper.ToGoogleCallbackResponse(response))
}

func (h *AuthHandler) SpotifyAuth(c echo.Context) error {
	var dto dtos.SpotifyAuthURLRequest
	if err := c.Bind(&dto); err != nil {
		dto.State = "random-state"
	}

	if dto.State == "" {
		dto.State = "random-state"
	}

	request := h.mapper.ToSpotifyAuthURLRequest(&dto)

	// Check if this is a link request (user is authenticated)
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" {
		// This is a link request, use link use case
		response, err := h.linkSpotifyUseCase.GetAuthURL(c.Request().Context(), *request)
		if err != nil {
			h.logger.Sugar().Errorf("Failed to generate Spotify auth URL for linking: %v", err)
			return HandleUseCaseError(c, err)
		}
		h.logger.Sugar().Infof("Generated Spotify auth URL for linking")

		// Check if client accepts JSON
		acceptHeader := c.Request().Header.Get("Accept")
		if acceptHeader == string(ContentTypeApplicationJson) || c.Request().Header.Get("Content-Type") == string(ContentTypeApplicationJson) {
			return SendSuccess(c, http.StatusOK, h.mapper.ToSpotifyAuthURLResponse(response))
		}
		return c.Redirect(http.StatusFound, response.URL)
	}

	// Regular login flow
	response, err := h.spotifyLoginUseCase.GetAuthURL(c.Request().Context(), *request)
	if err != nil {
		h.logger.Sugar().Errorf("Failed to generate Spotify auth URL: %v", err)
		return HandleUseCaseError(c, err)
	}

	h.logger.Sugar().Infof("Generated Spotify auth URL")

	// Check if client accepts JSON
	acceptHeader := c.Request().Header.Get("Accept")
	if acceptHeader == string(ContentTypeApplicationJson) || c.Request().Header.Get("Content-Type") == string(ContentTypeApplicationJson) {
		return SendSuccess(c, http.StatusOK, h.mapper.ToSpotifyAuthURLResponse(response))
	}
	return c.Redirect(http.StatusFound, response.URL)
}

func (h *AuthHandler) SpotifyCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		h.logger.Sugar().Warn("Missing authorization code in Spotify callback")
		return SendError(c, http.StatusBadRequest, "invalid_request", "Missing authorization code")
	}

	// Check if this is a link request (user is authenticated)
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" {
		// This is a link request - extract user from JWT
		claims, err := GetUserFromJWT(c)
		if err != nil {
			h.logger.Sugar().Warnf("Invalid JWT token in Spotify link callback: %v", err)
			return SendError(c, http.StatusUnauthorized, "unauthorized", "Invalid or missing authentication token")
		}

		dto := dtos.LinkSpotifyRequest{
			Code:  code,
			State: state,
		}

		request := h.mapper.ToLinkSpotifyRequest(&dto, claims.UserID.String())

		response, err := h.linkSpotifyUseCase.Execute(c.Request().Context(), *request)
		if err != nil {
			h.logger.Sugar().Errorf("Spotify link failed for user %s: %v", claims.UserID, err)
			return HandleUseCaseError(c, err)
		}

		h.logger.Sugar().Infof("Spotify account linked successfully for user: %s", claims.UserID)
		return SendSuccess(c, http.StatusOK, h.mapper.ToLinkSpotifyResponse(response))
	}

	// Regular login flow
	dto := dtos.SpotifyCallbackRequest{
		Code:  code,
		State: state,
	}

	request := h.mapper.ToSpotifyCallbackRequest(&dto)

	response, err := h.spotifyLoginUseCase.HandleCallback(c.Request().Context(), *request)
	if err != nil {
		h.logger.Sugar().Errorf("Spotify callback failed: %v", err)
		return HandleUseCaseError(c, err)
	}

	if response.IsNewUser {
		h.logger.Sugar().Infof("New user created via Spotify OAuth: %s", response.UserID)
	} else {
		h.logger.Sugar().Infof("User logged in via Spotify OAuth: %s", response.UserID)
	}

	return SendSuccess(c, http.StatusOK, h.mapper.ToSpotifyCallbackResponse(response))
}

func GetUserFromJWT(c echo.Context) (*middleware.Claims, error) {
	return middleware.GetUserFromContext(c)
}
