package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zandomed/sync-playlist-api/internal/usecases/health"
)

type HealthHandler struct {
	getStatusUseCase *health.GetStatusUseCase
}

func NewHealthHandler(getStatusUseCase *health.GetStatusUseCase) *HealthHandler {
	return &HealthHandler{
		getStatusUseCase: getStatusUseCase,
	}
}

func (h *HealthHandler) GetStatus(c echo.Context) error {
	status := h.getStatusUseCase.Execute()
	return c.JSON(http.StatusOK, status)
}