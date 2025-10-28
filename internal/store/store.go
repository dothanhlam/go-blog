package store

import "go-blog/internal/model"

// UserStore defines the interface for user data persistence.
type UserStore interface {
	Create(user *model.User) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	// Login is handled in the service layer, so no change here is needed.
	GetByID(id int) (*model.User, error)
}

// PostStore defines the interface for post data persistence.
type PostStore interface {
	Create(post *model.Post) (*model.Post, error)
	Update(post *model.Post) (*model.Post, error)
	GetByID(id int) (*model.Post, error)
	List(limit, offset int) ([]*model.Post, error)
	CreateHistory(history *model.PostHistory) error
	Search(query string, limit, offset int) ([]*model.Post, error)
}