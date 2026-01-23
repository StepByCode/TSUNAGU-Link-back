package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StepByCode/TSUNAGU-Link-back/internal/service"
	"github.com/StepByCode/TSUNAGU-Link-back/internal/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRespondWithError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := RespondWithError(c, http.StatusBadRequest, "test error")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "test error")
}

func TestRespondWithValidationError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	errs := validator.ValidationErrors{
		"email": validator.FieldError{
			Message: "must be a valid email address",
			Value:   "invalid",
		},
	}

	err := RespondWithValidationError(c, errs)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "validation failed")
	assert.Contains(t, rec.Body.String(), "must be a valid email address")
}

func TestHandleServiceError_UserAlreadyExists(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handleServiceError(c, service.ErrUserAlreadyExists)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
	assert.Contains(t, rec.Body.String(), "email already exists")
}

func TestHandleServiceError_UserNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handleServiceError(c, service.ErrUserNotFound)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "user not found")
}

func TestHandleServiceError_InvalidCredentials(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handleServiceError(c, service.ErrInvalidCredentials)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid credentials")
}

func TestHandleServiceError_InternalError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	unexpectedErr := errors.New("unexpected database error")
	err := handleServiceError(c, unexpectedErr)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "internal server error")
	// Note: The actual error is logged, but not returned to client
}
