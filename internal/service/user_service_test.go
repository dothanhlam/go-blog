package service_test

import (
	"errors"
	"go-blog/internal/model"
	"go-blog/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserStore is a mock implementation of store.UserStore
type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) Create(user *model.User) (*model.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) GetByID(id int) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) GetByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func TestUserService_Register(t *testing.T) {
	mockStore := new(MockUserStore)
	userService := service.NewUserService(mockStore)

	password := "password123"
	user := &model.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: password,
	}

	createdUser := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		// Password will be hashed, so we can't strict match it in the mock return for equality easily unless we capture it
		// But for the mock return value, we simulate what the DB returns
		Password: "hashed_password_from_db",
	}

	// Use Run to validate the input password was hashed
	mockStore.On("Create", mock.MatchedBy(func(u *model.User) bool {
		// Verify password was changed (hashed)
		return u.Password != password && u.Username == user.Username
	})).Return(createdUser, nil).Once()

	// Execute
	result, err := userService.Register(user)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdUser.ID, result.ID)
	assert.Equal(t, createdUser.Username, result.Username)

	mockStore.AssertExpectations(t)
}

func TestUserService_Register_Error(t *testing.T) {
	mockStore := new(MockUserStore)
	userService := service.NewUserService(mockStore)

	user := &model.User{
		Username: "testuser",
		Password: "password123",
	}

	expectedErr := errors.New("db error")

	mockStore.On("Create", mock.Anything).Return(nil, expectedErr).Once()

	// Execute
	result, err := userService.Register(user)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	mockStore.AssertExpectations(t)
}

func TestUserService_GetByID(t *testing.T) {
	mockStore := new(MockUserStore)
	userService := service.NewUserService(mockStore)

	userID := 1
	expectedUser := &model.User{
		ID:       userID,
		Username: "testuser",
	}

	mockStore.On("GetByID", userID).Return(expectedUser, nil).Once()

	// Execute
	result, err := userService.GetByID(userID)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)

	mockStore.AssertExpectations(t)
}

func TestUserService_GetByID_Error(t *testing.T) {
	mockStore := new(MockUserStore)
	userService := service.NewUserService(mockStore)

	userID := 1
	expectedErr := errors.New("not found")

	mockStore.On("GetByID", userID).Return(nil, expectedErr).Once()

	// Execute
	result, err := userService.GetByID(userID)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)

	mockStore.AssertExpectations(t)
}

func TestUserService_Login(t *testing.T) {
	mockStore := new(MockUserStore)
	userService := service.NewUserService(mockStore)

	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &model.User{
		ID:       1,
		Email:    email,
		Password: string(hashedPassword),
	}

	mockStore.On("GetByEmail", email).Return(user, nil).Once()

	// Execute
	result, err := userService.Login(email, password)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user, result)

	mockStore.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	mockStore := new(MockUserStore)
	userService := service.NewUserService(mockStore)

	email := "nonexistent@example.com"
	password := "password123"

	mockStore.On("GetByEmail", email).Return(nil, errors.New("not found")).Once()

	// Execute
	result, err := userService.Login(email, password)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, service.ErrNotFound, err)

	mockStore.AssertExpectations(t)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	mockStore := new(MockUserStore)
	userService := service.NewUserService(mockStore)

	email := "test@example.com"
	password := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	user := &model.User{
		ID:       1,
		Email:    email,
		Password: string(hashedPassword),
	}

	mockStore.On("GetByEmail", email).Return(user, nil).Once()

	// Execute
	result, err := userService.Login(email, password)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, service.ErrPermissionDenied, err)

	mockStore.AssertExpectations(t)
}
