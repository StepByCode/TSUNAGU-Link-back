package handler

import (
	"net/http"
	"strconv"

	"github.com/StepByCode/TSUNAGU-Link-back/internal/model"
	"github.com/StepByCode/TSUNAGU-Link-back/internal/service"
	"github.com/StepByCode/TSUNAGU-Link-back/internal/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
	validator   *validator.Validator
}

func NewUserHandler(userService service.UserService, validator *validator.Validator) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req model.CreateUserRequest

	// JSON unmarshalling
	if err := c.Bind(&req); err != nil {
		return RespondWithError(c, http.StatusBadRequest, "invalid request format")
	}

	// Validation
	if errs := h.validator.Validate(&req); errs.HasErrors() {
		return RespondWithValidationError(c, errs)
	}

	// Business logic
	user, err := h.userService.CreateUser(c.Request().Context(), &req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return RespondWithError(c, http.StatusBadRequest, "invalid user id")
	}

	user, err := h.userService.GetUser(c.Request().Context(), id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return RespondWithError(c, http.StatusBadRequest, "invalid user id")
	}

	var req model.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return RespondWithError(c, http.StatusBadRequest, "invalid request format")
	}

	// Validation
	if errs := h.validator.Validate(&req); errs.HasErrors() {
		return RespondWithValidationError(c, errs)
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), id, &req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return RespondWithError(c, http.StatusBadRequest, "invalid user id")
	}

	if err := h.userService.DeleteUser(c.Request().Context(), id); err != nil {
		return handleServiceError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}

	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if offset < 0 {
		offset = 0
	}

	users, err := h.userService.ListUsers(c.Request().Context(), limit, offset)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return RespondWithError(c, http.StatusBadRequest, "invalid request format")
	}

	// Validation
	if errs := h.validator.Validate(&req); errs.HasErrors() {
		return RespondWithValidationError(c, errs)
	}

	response, err := h.userService.Login(c.Request().Context(), &req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	api := e.Group("/api/v1")

	api.POST("/auth/login", h.Login)
	api.POST("/users", h.CreateUser)
	api.GET("/users/:id", h.GetUser)
	api.PUT("/users/:id", h.UpdateUser)
	api.DELETE("/users/:id", h.DeleteUser)
	api.GET("/users", h.ListUsers)
}
