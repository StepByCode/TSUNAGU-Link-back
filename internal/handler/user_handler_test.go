package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/StepByCode/TSUNAGU-Link-back/internal/model"
	"github.com/StepByCode/TSUNAGU-Link-back/internal/service"
	"github.com/StepByCode/TSUNAGU-Link-back/internal/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id uuid.UUID, req *model.UpdateUserRequest) (*model.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LoginResponse), args.Error(1)
}

func TestUserHandler_CreateUser(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	reqBody := `{"email":"test@example.com","name":"Test User","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedUser := &model.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.CreateUserRequest")).Return(expectedUser, nil)

	err := handler.CreateUser(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var user model.User
	err = json.Unmarshal(rec.Body.Bytes(), &user)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Name, user.Name)

	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_InvalidRequest(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	reqBody := `{"invalid json`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateUser(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_CreateUser_ValidationError_InvalidEmail(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	reqBody := `{"email":"invalid-email","name":"Test User","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateUser(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "validation failed")
	assert.Contains(t, rec.Body.String(), "email")
}

func TestUserHandler_CreateUser_ValidationError_ShortPassword(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	reqBody := `{"email":"test@example.com","name":"Test User","password":"short"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateUser(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "validation failed")
	assert.Contains(t, rec.Body.String(), "password")
}

func TestUserHandler_CreateUser_ValidationError_MissingFields(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	reqBody := `{"email":"","name":"","password":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateUser(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "validation failed")
}

func TestUserHandler_CreateUser_ServiceError(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	reqBody := `{"email":"test@example.com","name":"Test User","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.CreateUserRequest")).Return(nil, fmt.Errorf("service error"))

	err := handler.CreateUser(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	userID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	expectedUser := &model.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("GetUser", mock.Anything, userID).Return(expectedUser, nil)

	err := handler.GetUser(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var user model.User
	err = json.Unmarshal(rec.Body.Bytes(), &user)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)

	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_InvalidID(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/invalid-id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid-id")

	err := handler.GetUser(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	userID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	// Use the actual domain error
	mockService.On("GetUser", mock.Anything, userID).Return(nil, service.ErrUserNotFound)

	err := handler.GetUser(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateUser(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	userID := uuid.New()
	reqBody := `{"name":"Updated Name","email":"updated@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+userID.String(), strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	updatedUser := &model.User{
		ID:        userID,
		Email:     "updated@example.com",
		Name:      "Updated Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("UpdateUser", mock.Anything, userID, mock.AnythingOfType("*model.UpdateUserRequest")).Return(updatedUser, nil)

	err := handler.UpdateUser(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var user model.User
	err = json.Unmarshal(rec.Body.Bytes(), &user)
	require.NoError(t, err)
	assert.Equal(t, updatedUser.Name, user.Name)
	assert.Equal(t, updatedUser.Email, user.Email)

	mockService.AssertExpectations(t)
}

func TestUserHandler_DeleteUser(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	userID := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+userID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	mockService.On("DeleteUser", mock.Anything, userID).Return(nil)

	err := handler.DeleteUser(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)

	mockService.AssertExpectations(t)
}

func TestUserHandler_ListUsers(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?limit=10&offset=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedUsers := []*model.User{
		{
			ID:        uuid.New(),
			Email:     "user1@example.com",
			Name:      "User 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Email:     "user2@example.com",
			Name:      "User 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockService.On("ListUsers", mock.Anything, 10, 0).Return(expectedUsers, nil)

	err := handler.ListUsers(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var users []*model.User
	err = json.Unmarshal(rec.Body.Bytes(), &users)
	require.NoError(t, err)
	assert.Len(t, users, 2)

	mockService.AssertExpectations(t)
}

func TestUserHandler_ListUsers_DefaultPagination(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("ListUsers", mock.Anything, 10, 0).Return([]*model.User{}, nil)

	err := handler.ListUsers(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockService.AssertExpectations(t)
}

func TestUserHandler_Login(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	reqBody := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedResponse := &model.LoginResponse{
		Token: "test-token",
		User: model.User{
			ID:        uuid.New(),
			Email:     "test@example.com",
			Name:      "Test User",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockService.On("Login", mock.Anything, mock.AnythingOfType("*model.LoginRequest")).Return(expectedResponse, nil)

	err := handler.Login(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response model.LoginResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse.Token, response.Token)
	assert.Equal(t, expectedResponse.User.Email, response.User.Email)

	mockService.AssertExpectations(t)
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := new(MockUserService)
	v := validator.NewValidator()
	handler := NewUserHandler(mockService, v)

	e := echo.New()
	reqBody := `{"email":"test@example.com","password":"wrongpassword"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.On("Login", mock.Anything, mock.AnythingOfType("*model.LoginRequest")).Return(nil, service.ErrInvalidCredentials)

	err := handler.Login(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	mockService.AssertExpectations(t)
}
