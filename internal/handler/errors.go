package handler

import (
	"errors"
	"net/http"

	"github.com/StepByCode/TSUNAGU-Link-back/internal/service"
	"github.com/StepByCode/TSUNAGU-Link-back/internal/validator"
	"github.com/labstack/echo/v4"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondWithError sends a standardized error response
func RespondWithError(c echo.Context, code int, message string) error {
	return c.JSON(code, ErrorResponse{Error: message})
}

// RespondWithValidationError sends validation error details
func RespondWithValidationError(c echo.Context, errs validator.ValidationErrors) error {
	return c.JSON(http.StatusBadRequest, errs.ToResponse())
}

// handleServiceError maps service layer errors to appropriate HTTP status codes
func handleServiceError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrUserAlreadyExists):
		return RespondWithError(c, http.StatusConflict, "email already exists")
	case errors.Is(err, service.ErrUserNotFound):
		return RespondWithError(c, http.StatusNotFound, "user not found")
	case errors.Is(err, service.ErrInvalidCredentials):
		return RespondWithError(c, http.StatusUnauthorized, "invalid credentials")
	default:
		// Log the actual error for debugging, but return generic message to client
		c.Logger().Errorf("Internal server error: %v, request: %s %s", err, c.Request().Method, c.Request().URL.Path)
		return RespondWithError(c, http.StatusInternalServerError, "internal server error")
	}
}
