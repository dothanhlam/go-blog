package api

import (
	"go-blog/internal/service"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/labstack/echo/v4"
)

// WebHandler handles requests for server-side rendered pages.
type WebHandler struct {
	postService service.PostService
}

// NewWebHandler creates a new WebHandler.
func NewWebHandler(ps service.PostService) *WebHandler {
	return &WebHandler{postService: ps}
}

// RenderIndexPage renders the home page with a list of all posts.
func (h *WebHandler) RenderIndexPage(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10 // Default limit
	}

	posts, err := h.postService.List(page, limit)
	if err != nil {
		// In a real app, you'd want to log this error
		return c.String(http.StatusInternalServerError, "Could not fetch posts")
	}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Posts": posts,
	})
}

// RenderPostPage renders the page for a single post.
func (h *WebHandler) RenderPostPage(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid post ID")
	}

	post, mdContent, err := h.postService.GetByID(id)
	if err != nil {
		// Assuming service returns a specific error for not found
		return c.String(http.StatusNotFound, "Post not found")
	}

	// Convert markdown to HTML
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	htmlContent := markdown.ToHTML([]byte(mdContent), p, nil)

	return c.Render(http.StatusOK, "post.html", map[string]interface{}{
		"Title":   post.Title,
		"Content": template.HTML(htmlContent), // Use template.HTML to prevent escaping
	})
}