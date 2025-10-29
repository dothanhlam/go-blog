package service

import (
	"fmt"
	
	"go-blog/internal/store"
	"go-blog/internal/model"
	"go-blog/internal/storage"
)

type PostService interface {
	Create(title, content string, userID int) (*model.Post, error)
	GetByID(id int) (*model.Post, string, error)
	List(page, limit int) ([]*model.Post, error)
	Update(postID int, title, content string, userID int) (*model.Post, error)
	Search(query string, page, limit int) ([]*model.Post, error)
}

type postService struct {
	postStore   store.PostStore
	fileStorage storage.FileStorage
}

func NewPostService(ps store.PostStore, fs storage.FileStorage) PostService {
	return &postService{postStore: ps, fileStorage: fs}
}

func (s *postService) Create(title, content string, userID int) (*model.Post, error) {
	// Create the post metadata in the database first to get an ID.
	post := &model.Post{
		UserID:  userID,
		Title:   title,
		Version: 1,
	}

	createdPost, err := s.postStore.Create(post)
	if err != nil {
		return nil, err
	}

	// Use the post ID to create a unique path for the content file.
	// Path format: user_<userID>/post_<postID>_v1.md
	contentPath := fmt.Sprintf("user_%d/post_%d_v%d.md", userID, createdPost.ID, createdPost.Version)

	// Save the markdown content to the configured storage (local or S3).
	if err := s.fileStorage.Save(contentPath, []byte(content)); err != nil {
		// TODO: Consider rolling back the database transaction if file storage fails.
		return nil, err
	}

	// Update the post record with the content path.
	createdPost.ContentPath = contentPath
	return s.postStore.Update(createdPost)
}

func (s *postService) GetByID(id int) (*model.Post, string, error) {
	post, err := s.postStore.GetByID(id)
	if err != nil {
		return nil, "", ErrNotFound
	}

	content, err := s.fileStorage.Read(post.ContentPath)
	if err != nil {
		// If we can't read the file, the post is in an inconsistent state.
		return nil, "", fmt.Errorf("could not read content for post %d: %w", id, err)
	}

	return post, string(content), nil
}

func (s *postService) List(page, limit int) ([]*model.Post, error) {
	offset := (page - 1) * limit
	return s.postStore.List(limit, offset)
}

func (s *postService) Search(query string, page, limit int) ([]*model.Post, error) {
	offset := (page - 1) * limit
	return s.postStore.Search(query, limit, offset)
}

func (s *postService) Update(postID int, title, content string, userID int) (*model.Post, error) {
	// TODO: This entire operation should be in a single database transaction.

	// 1. Get the current post from the database.
	post, err := s.postStore.GetByID(postID)
	if err != nil {
		return nil, ErrNotFound
	}

	// 2. Verify ownership.
	if post.UserID != userID {
		return nil, ErrPermissionDenied
	}

	// 3. Create a history record for the *current* version before we update it.
	history := &model.PostHistory{
		PostID:      post.ID,
		Version:     post.Version,
		ContentPath: post.ContentPath,
	}
	if err := s.postStore.CreateHistory(history); err != nil {
		return nil, fmt.Errorf("failed to create post history: %w", err)
	}

	// 4. Increment version and define the new content path.
	newVersion := post.Version + 1
	newContentPath := fmt.Sprintf("user_%d/post_%d_v%d.md", post.UserID, post.ID, newVersion)

	// 5. Save the new content to file storage.
	if err := s.fileStorage.Save(newContentPath, []byte(content)); err != nil {
		// If this fails, we have a history record but haven't updated the main post.
		// A transaction would allow us to roll back the history creation.
		return nil, fmt.Errorf("failed to save new post content: %w", err)
	}

	// 6. Update the post model with new data.
	post.Title = title
	post.Version = newVersion
	post.ContentPath = newContentPath

	// 7. Persist the updated post to the database.
	return s.postStore.Update(post)
}