package api

import (
	"errors"
	"go-blog/internal/config"
	"go-blog/internal/middleware" // Added this import
	"go-blog/internal/service"
	"html/template"
	"net/http"
	"strconv"
	"time"
	"log" // Added log import

	"github.com/golang-jwt/jwt/v5" // Ensured this import is present
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/labstack/echo/v4"
)

// WebHandler handles requests for server-side rendered pages.
type WebHandler struct {
	cfg         *config.Config
	postService service.PostService
	userService service.UserService
}

// NewWebHandler creates a new WebHandler.
func NewWebHandler(cfg *config.Config, ps service.PostService, us service.UserService) *WebHandler {
	return &WebHandler{cfg: cfg, postService: ps, userService: us}
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
		log.Printf("Error fetching posts: %v", err)
		return c.String(http.StatusInternalServerError, "Could not fetch posts")
	}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"User":    c.Get(middleware.UserContextKey),
		"Context": c,
		"Posts":   posts,
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
		if errors.Is(err, service.ErrNotFound) {
			// Use the localizer to return a translated "not found" message.
			// This requires a bit more setup in the handler or a custom error handler.
			// For now, we keep it simple.
			return c.String(http.StatusNotFound, "post_not_found")
		}
		return c.String(http.StatusNotFound, "post_not_found")
	}

	// Convert markdown to HTML
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	htmlContent := markdown.ToHTML([]byte(mdContent), p, nil)
	// Use html.Renderer to get more control if needed in the future.
	// renderer := html.NewRenderer(html.RendererOptions{Flags: html.CommonFlags})
	// htmlContent = markdown.ToHTML([]byte(mdContent), p, renderer)
	
	// log.Printf("web handler ", c.Get(middleware.UserContextKey))
	
	return c.Render(http.StatusOK, "post.html", map[string]interface{}{
		"User":    c.Get(middleware.UserContextKey),
		"Context": c,
		"Post":    post,
		"Content": template.HTML(htmlContent), // Use template.HTML to prevent escaping
	})
}

// RenderLoginPage renders the login page.
func (h *WebHandler) RenderLoginPage(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

// HandleLogin processes the login form submission.
func (h *WebHandler) HandleLogin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := h.userService.Login(email, password)
	if err != nil {
		// In a real app, you'd render the login page again with an error message.
		return c.Redirect(http.StatusFound, "/login?error=invalid_credentials")
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	expirationTime := time.Now().Add(time.Hour * time.Duration(h.cfg.TokenExpiresInHours))
	claims["id"] = user.ID
	claims["exp"] = expirationTime.Unix()

	// Generate encoded token
	t, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return c.Redirect(http.StatusFound, "/login?error=token_error")
	}

	// Set the token in a secure, http-only cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    t,
		Expires:  expirationTime,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to false for local HTTP development, true for HTTPS in production
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}

// HandleLogout clears the session cookie and logs the user out.
func (h *WebHandler) HandleLogout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Expire immediately
		Path:     "/",
		Secure:   false, // Consistent with login for local development
		HttpOnly: true,
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, "/")
}

// CustomHTTPErrorHandler handles HTTP errors for the web interface.
func (h *WebHandler) CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	// If it's a 404, render our custom 404 page
	if code == http.StatusNotFound {
		// We can pass the user from the context to the template,
		// so the header/navigation still looks correct.
		err := c.Render(http.StatusNotFound, "404.html", map[string]interface{}{
			"Context": c,
			"User": c.Get(middleware.UserContextKey),
		})
		if err != nil {
			c.Logger().Error(err)
		}
		return
	}
}