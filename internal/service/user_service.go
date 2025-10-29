package service

import (
	"go-blog/internal/model"
	"go-blog/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user-related business logic.
type UserService interface {
	Login(email, password string) (*model.User, error)
	GetByID(id int) (*model.User, error)
	Register(user *model.User) (*model.User, error)
}

type userService struct {
	userStore store.UserStore
}

// NewUserService creates a new UserService.
func NewUserService(us store.UserStore) UserService {
	return &userService{userStore: us}
}

// GetByID retrieves a user by their ID.
func (s *userService) GetByID(id int) (*model.User, error) {
	return s.userStore.GetByID(id)
}

// Register creates a new user after hashing their password.
func (s *userService) Register(user *model.User) (*model.User, error) {
	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	return s.userStore.Create(user)
}

// Login authenticates a user and returns the user model if successful.
func (s *userService) Login(email, password string) (*model.User, error) {
	user, err := s.userStore.GetByEmail(email)
	if err != nil {
		return nil, ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrPermissionDenied
	}

	return user, nil
}