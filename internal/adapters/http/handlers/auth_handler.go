package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zandomed/sync-playlist-api/internal/adapters/http/dtos"
	"github.com/zandomed/sync-playlist-api/internal/adapters/http/mappers"
	"github.com/zandomed/sync-playlist-api/internal/usecases/auth"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

type AuthHandler struct {
	registerUseCase *auth.RegisterUserUseCase
	loginUseCase    *auth.LoginUserUseCase
	mapper          *mappers.AuthMapper
	logger          *logger.Logger
}

func NewAuthHandler(
	registerUseCase *auth.RegisterUserUseCase,
	loginUseCase *auth.LoginUserUseCase,
	mapper *mappers.AuthMapper,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		registerUseCase: registerUseCase,
		loginUseCase:    loginUseCase,
		mapper:          mapper,
		logger:          logger,
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
