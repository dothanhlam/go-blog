package service

import (
	"go-blog/internal/model"
	"go-blog/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *model.User) (*model.User, error)
	Login(email, password string) (*model.User, error)
}

type userService struct {
	userStore store.UserStore
}

func NewUserService(us store.UserStore) UserService {
	return &userService{userStore: us}
}

func (s *userService) Register(user *model.User) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	return s.userStore.Create(user)
}

func (s *userService) Login(email, password string) (*model.User, error) {
	// 1. Retrieve user by email from the database
	user, err := s.userStore.GetByEmail(email)
	if err != nil {
		// User not found
		return nil, err
	}

	// 2. Compare the provided password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// Passwords don't match
		return nil, err
	}

	return user, nil
}