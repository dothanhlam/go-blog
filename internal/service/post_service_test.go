package service_test

import (
	"fmt"
	"go-blog/internal/model"
	"go-blog/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostStore is a mock implementation of store.PostStore
type MockPostStore struct {
	mock.Mock
}

func (m *MockPostStore) Create(post *model.Post) (*model.Post, error) {
	args := m.Called(post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostStore) Update(post *model.Post) (*model.Post, error) {
	args := m.Called(post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostStore) GetByID(id int) (*model.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostStore) List(limit, offset int) ([]*model.Post, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockPostStore) CreateHistory(history *model.PostHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockPostStore) Search(query string, limit, offset int) ([]*model.Post, error) {
	args := m.Called(query, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Post), args.Error(1)
}

// MockFileStorage is a mock implementation of storage.FileStorage
type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) Save(path string, data []byte) error {
	args := m.Called(path, data)
	return args.Error(0)
}

func (m *MockFileStorage) Read(path string) ([]byte, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func TestPostService_Create(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockFileStorage := new(MockFileStorage)
	postSvc := service.NewPostService(mockPostStore, mockFileStorage)

	userID := 1
	title := "Test Title"
	content := "Test Content"

	subTitle := "Test SubTitle"
	image := "test_image.jpg"
	tags := []string{"tag1", "tag2"}

	// The initial post model passed to Create
	initialPost := &model.Post{
		UserID:   userID,
		Title:    title,
		SubTitle: subTitle,
		Image:    image,
		Tags:     tags,
		Version:  1,
	}

	// The post model returned by the first Create call, now with an ID
	createdPostWithID := &model.Post{
		ID:       1,
		UserID:   userID,
		Title:    title,
		SubTitle: subTitle,
		Image:    image,
		Tags:     tags,
		Version:  1,
	}

	contentPath := fmt.Sprintf("user_%d/post_%d_v%d.md", userID, createdPostWithID.ID, createdPostWithID.Version)

	// The final post model after being updated with the content path
	finalPost := &model.Post{
		ID:          1,
		UserID:      userID,
		Title:       title,
		SubTitle:    subTitle,
		Image:       image,
		Tags:        tags,
		Version:     1,
		ContentPath: contentPath,
	}

	// Setup mock expectations
	mockPostStore.On("Create", initialPost).Return(createdPostWithID, nil).Once()
	mockFileStorage.On("Save", contentPath, []byte(content)).Return(nil).Once()
	mockPostStore.On("Update", finalPost).Return(finalPost, nil).Once()

	// Execute the service method
	// Execute the service method
	post, err := postSvc.Create(title, subTitle, image, tags, content, userID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, finalPost, post)

	// Verify that all expectations were met
	mockPostStore.AssertExpectations(t)
	mockFileStorage.AssertExpectations(t)
}

func TestPostService_GetByID(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockFileStorage := new(MockFileStorage)
	postSvc := service.NewPostService(mockPostStore, mockFileStorage)

	postID := 1
	contentPath := "user_1/post_1_v1.md"
	content := "This is the content."

	dbPost := &model.Post{
		ID:          postID,
		UserID:      1,
		Title:       "A Post",
		ContentPath: contentPath,
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Setup mock expectations
	mockPostStore.On("GetByID", postID).Return(dbPost, nil).Once()
	mockFileStorage.On("Read", contentPath).Return([]byte(content), nil).Once()

	// Execute
	post, postContent, err := postSvc.GetByID(postID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, dbPost, post)
	assert.Equal(t, content, postContent)

	mockPostStore.AssertExpectations(t)
	mockFileStorage.AssertExpectations(t)
}

func TestPostService_Update(t *testing.T) {
	mockPostStore := new(MockPostStore)
	mockFileStorage := new(MockFileStorage)
	postSvc := service.NewPostService(mockPostStore, mockFileStorage)

	postID := 1
	userID := 1
	originalVersion := 1
	originalContentPath := fmt.Sprintf("user_%d/post_%d_v%d.md", userID, postID, originalVersion)

	currentPost := &model.Post{
		ID:          postID,
		UserID:      userID,
		Title:       "Original Title",
		ContentPath: originalContentPath,
		Version:     originalVersion,
	}

	newTitle := "Updated Title"
	newSubTitle := "Updated SubTitle"
	newImage := "updated_image.jpg"
	newTags := []string{"updated_tag1", "updated_tag2"}
	newContent := "Updated Content"
	newVersion := 2
	newContentPath := fmt.Sprintf("user_%d/post_%d_v%d.md", userID, postID, newVersion)

	// Setup mock expectations
	mockPostStore.On("GetByID", postID).Return(currentPost, nil).Once()
	mockPostStore.On("CreateHistory", mock.AnythingOfType("*model.PostHistory")).Return(nil).Once()
	mockFileStorage.On("Save", newContentPath, []byte(newContent)).Return(nil).Once()
	mockPostStore.On("Update", mock.AnythingOfType("*model.Post")).Return(currentPost, nil).Once()

	// Execute
	// Execute
	updatedPost, err := postSvc.Update(postID, newTitle, newSubTitle, newImage, newTags, newContent, userID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, updatedPost)
	assert.Equal(t, newTitle, updatedPost.Title)
	assert.Equal(t, newSubTitle, updatedPost.SubTitle)
	assert.Equal(t, newImage, updatedPost.Image)
	assert.Equal(t, newTags, updatedPost.Tags)
	assert.Equal(t, newVersion, updatedPost.Version)
	assert.Equal(t, newContentPath, updatedPost.ContentPath)

	mockPostStore.AssertExpectations(t)
	mockFileStorage.AssertExpectations(t)
}