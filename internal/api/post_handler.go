package api

import (
	"go-blog/internal/service"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(ps service.PostService) *PostHandler {
	return &PostHandler{postService: ps}
}

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"` // Markdown content
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	// Get user ID from JWT token claims
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := int(claims["id"].(float64))

	var req CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	post, err := h.postService.Create(req.Title, req.Content, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) UpdatePost(c echo.Context) error {
	// Get user ID from JWT token claims for ownership verification
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := int(claims["id"].(float64))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	var req UpdatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	post, err := h.postService.Update(id, req.Title, req.Content, userID)
	if err != nil {
		// This could be a not found error, a permission error, or a server error
		// A more robust error handling mechanism would be better here.
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, post)
}

func (h *PostHandler) GetPost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
	}

	post, content, err := h.postService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"post":    post,
		"content": content,
	})
}

func (h *PostHandler) ListPosts(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	posts, err := h.postService.List(page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, posts)
}