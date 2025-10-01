package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/infra/http/dtos"
)

type ContentTypeAccept string

const (
	ContentTypeApplicationJson ContentTypeAccept = "application/json"
)

func SendSuccess(c echo.Context, status int, data interface{}) error {
	response := dtos.SuccessResponse{
		Data: data,
	}
	return c.JSON(status, response)
}

func SendError(c echo.Context, status int, code, message string) error {
	response := dtos.ErrorResponse{
		Error:   code,
		Message: message,
		Code:    code,
	}
	return c.JSON(status, response)
}

func SendValidationError(c echo.Context, err error) error {
	response := dtos.ErrorResponse{
		Error:   "validation_failed",
		Message: "Request validation failed",
		Code:    "validation_failed",
	}
	return c.JSON(http.StatusBadRequest, response)
}

func HandleUseCaseError(c echo.Context, err error) error {
	switch e := err.(type) {
	case *errors.DomainError:
		return SendError(c, http.StatusBadRequest, e.Code(), e.Message())
	case *errors.AuthenticationError:
		return SendError(c, http.StatusUnauthorized, e.Code(), e.Message())
	case *errors.ValidationError:
		return SendError(c, http.StatusBadRequest, e.Code(), e.Message())
	case *errors.NotFoundError:
		return SendError(c, http.StatusNotFound, e.Code(), e.Message())
	default:
		return SendError(c, http.StatusInternalServerError, "internal_error", "An internal error occurred")
	}
}
