package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/services"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

type AuthHandler struct {
	service services.AuthService
	logger  *logger.Logger
	cfg     *config.Config
	// Add fields as necessary, e.g., services, logger, etc.
}

func NewAuthHandler(s services.AuthService, cfg *config.Config, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		service: s,
		logger:  logger,
		cfg:     cfg,
	}
}

func (h *AuthHandler) LoginWithPass(c echo.Context) error {
	h.logger.Sugar().Info("LoginWithPass called")
	req := new(services.LoginRequest)
	if err := c.Bind(req); err != nil {
		return sendError(c, http.StatusBadRequest, err, "Invalid request body")
	}
	h.logger.Sugar().Infof("Login request: %+v", req)
	if err := c.Validate(req); err != nil {
		return sendValidationError(c, http.StatusBadRequest, err)
	}
	if resp, err := h.service.LoginWithPass(*req); err != nil {
		return sendError(c, http.StatusInternalServerError, err, "Failed to login user")
	} else {
		h.logger.Sugar().Infof("User logged in successfully: %s", req.Email)
		return c.JSON(http.StatusOK, resp)
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	// Implementation goes here
	h.logger.Sugar().Info("Register called")

	req := new(services.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return sendError(c, http.StatusBadRequest, err, "Invalid request body")
	}

	h.logger.Sugar().Infof("Register request: %+v", req)

	if err := c.Validate(req); err != nil {
		return sendValidationError(c, http.StatusBadRequest, err)
	}

	if resp, err := h.service.Register(*req); err != nil {
		return sendError(c, http.StatusInternalServerError, err, "Failed to register user")
	} else {
		h.logger.Sugar().Infof("User registered successfully: %d", resp.ID)
		return c.JSON(http.StatusCreated, resp)
	}
}
