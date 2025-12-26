package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/StepByCode/TSUNAGU-Link-back/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.User), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, "test-secret", 24)

	ctx := context.Background()
	req := &model.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	user, err := service.CreateUser(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Email, user.Email)
	assert.Equal(t, req.Name, user.Name)
	assert.NotEqual(t, req.Password, user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	assert.NoError(t, err, "password should be hashed correctly")

	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, "test-secret", 24)

	ctx := context.Background()
	req := &model.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(fmt.Errorf("db error"))

	user, err := service.CreateUser(ctx, req)

	require.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to create user")

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, "test-secret", 24)

	ctx := context.Background()
	userID := uuid.New()
	expectedUser := &model.User{
		ID:        userID,
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", ctx, userID).Return(expectedUser, nil)

	user, err := service.GetUser(ctx, userID)

	require.NoError(t, err)
	assert.Equal(t, expectedUser, user)

	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, "test-secret", 24)

	ctx := context.Background()
	userID := uuid.New()
	existingUser := &model.User{
		ID:        userID,
		Email:     "old@example.com",
		Name:      "Old Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newName := "New Name"
	newEmail := "new@example.com"
	req := &model.UpdateUserRequest{
		Name:  &newName,
		Email: &newEmail,
	}

	mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*model.User")).Return(nil)

	user, err := service.UpdateUser(ctx, userID, req)

	require.NoError(t, err)
	assert.Equal(t, newName, user.Name)
	assert.Equal(t, newEmail, user.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, "test-secret", 24)

	ctx := context.Background()
	userID := uuid.New()
	newName := "New Name"
	req := &model.UpdateUserRequest{
		Name: &newName,
	}

	mockRepo.On("GetByID", ctx, userID).Return(nil, fmt.Errorf("user not found"))

	user, err := service.UpdateUser(ctx, userID, req)

	require.Error(t, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, "test-secret", 24)

	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("Delete", ctx, userID).Return(nil)

	err := service.DeleteUser(ctx, userID)

	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestUserService_ListUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, "test-secret", 24)

	ctx := context.Background()
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

	mockRepo.On("List", ctx, 10, 0).Return(expectedUsers, nil)

	users, err := service.ListUsers(ctx, 10, 0)

	require.NoError(t, err)
	assert.Equal(t, expectedUsers, users)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, "test-secret", 24)

	ctx := context.Background()
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &model.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := &model.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)

	response, err := service.Login(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, user.ID, response.User.ID)
	assert.Equal(t, user.Email, response.User.Email)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, "test-secret", 24)

	ctx := context.Background()
	req := &model.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(nil, fmt.Errorf("user not found"))

	response, err := service.Login(ctx, req)

	require.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, "test-secret", 24)

	ctx := context.Background()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	user := &model.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := &model.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)

	response, err := service.Login(ctx, req)

	require.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockRepo.AssertExpectations(t)
}

func TestUserService_GenerateJWT(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo, "test-secret", 24).(*userService)

	user := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Name:  "Test User",
	}

	token, err := service.generateJWT(user)

	require.NoError(t, err)
	assert.NotEmpty(t, token)
}
