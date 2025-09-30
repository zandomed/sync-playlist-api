package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zandomed/sync-playlist-api/internal/adapters/http/dtos"
	"github.com/zandomed/sync-playlist-api/internal/adapters/http/mappers"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
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
		return h.sendError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := c.Validate(&dto); err != nil {
		h.logger.Sugar().Warnf("Validation failed: %v", err)
		return h.sendValidationError(c, err)
	}

	request := h.mapper.ToRegisterUserRequest(&dto)

	response, err := h.registerUseCase.Execute(c.Request().Context(), *request)
	if err != nil {
		h.logger.Sugar().Errorf("Registration failed: %v", err)
		return h.handleUseCaseError(c, err)
	}

	h.logger.Sugar().Infof("User registered successfully: %s", response.UserID)
	return h.sendSuccess(c, http.StatusCreated, h.mapper.ToRegisterResponse(response))
}

func (h *AuthHandler) Login(c echo.Context) error {
	var dto dtos.LoginRequest
	if err := c.Bind(&dto); err != nil {
		h.logger.Sugar().Warnf("Invalid request body: %v", err)
		return h.sendError(c, http.StatusBadRequest, "invalid_request", "Invalid request body")
	}

	if err := c.Validate(&dto); err != nil {
		h.logger.Sugar().Warnf("Validation failed: %v", err)
		return h.sendValidationError(c, err)
	}

	request := h.mapper.ToLoginUserRequest(&dto)

	response, err := h.loginUseCase.Execute(c.Request().Context(), *request)
	if err != nil {
		h.logger.Sugar().Warnf("Login failed for user %s: %v", dto.Email, err)
		return h.handleUseCaseError(c, err)
	}

	h.logger.Sugar().Infof("User logged in successfully: %s", response.UserID)
	return h.sendSuccess(c, http.StatusOK, h.mapper.ToLoginResponse(response))
}

func (h *AuthHandler) handleUseCaseError(c echo.Context, err error) error {
	switch e := err.(type) {
	case *errors.DomainError:
		return h.sendError(c, http.StatusBadRequest, e.Code(), e.Message())
	case *errors.AuthenticationError:
		return h.sendError(c, http.StatusUnauthorized, e.Code(), e.Message())
	case *errors.ValidationError:
		return h.sendError(c, http.StatusBadRequest, e.Code(), e.Message())
	case *errors.NotFoundError:
		return h.sendError(c, http.StatusNotFound, e.Code(), e.Message())
	default:
		return h.sendError(c, http.StatusInternalServerError, "internal_error", "An internal error occurred")
	}
}

func (h *AuthHandler) sendError(c echo.Context, status int, code, message string) error {
	response := dtos.ErrorResponse{
		Error:   code,
		Message: message,
		Code:    code,
	}
	return c.JSON(status, response)
}

func (h *AuthHandler) sendValidationError(c echo.Context, err error) error {
	response := dtos.ErrorResponse{
		Error:   "validation_failed",
		Message: "Request validation failed",
		Code:    "validation_failed",
	}
	return c.JSON(http.StatusBadRequest, response)
}

func (h *AuthHandler) sendSuccess(c echo.Context, status int, data interface{}) error {
	response := dtos.SuccessResponse{
		Data: data,
	}
	return c.JSON(status, response)
}
