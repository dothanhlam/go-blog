package mocks

import (
	"go-blog/internal/model"

	"github.com/stretchr/testify/mock"
)

// PostService is a mock for the service.PostService interface.
type PostService struct {
	mock.Mock
}

func (m *PostService) Create(title, content string, userID int) (*model.Post, error) {
	args := m.Called(title, content, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *PostService) Update(id int, title, content string, userID int) (*model.Post, error) {
	args := m.Called(id, title, content, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *PostService) GetByID(id int) (*model.Post, string, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*model.Post), args.String(1), args.Error(2)
}

func (m *PostService) List(page, limit int) ([]*model.Post, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *PostService) Search(query string, page, limit int) ([]*model.Post, error) {
	args := m.Called(query, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *PostService) Delete(id int, userID int) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *PostService) GetHistory(postID int) ([]*model.PostHistory, error) {
	args := m.Called(postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.PostHistory), args.Error(1)
}