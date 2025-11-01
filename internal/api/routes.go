package api

import (
	"go-blog/internal/config"
	"go-blog/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo-jwt/v4"
)

// RegisterRoutes sets up all the routes for the application.
func RegisterRoutes(e *echo.Echo, userService service.UserService, postService service.PostService, cfg *config.Config) {
	userHandler := NewUserHandler(userService)
	postHandler := NewPostHandler(postService)

	// API group
	apiGroup := e.Group("/api")

	// User routes
	apiGroup.POST("/register", userHandler.Register)
	apiGroup.POST("/login", userHandler.Login(cfg))

	// Post routes
	apiGroup.GET("/posts", postHandler.ListPosts) // Publicly accessible list of posts
	apiGroup.GET("/posts/:id", postHandler.GetPost)
	apiGroup.GET("/posts/search", postHandler.SearchPosts)

	// Authenticated routes
	authGroup := apiGroup.Group("")
	authGroup.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JWTSecret),
	}))
	authGroup.POST("/posts", postHandler.CreatePost)
	authGroup.PUT("/posts/:id", postHandler.UpdatePost)
	authGroup.POST("/posts/upload", postHandler.CreateFromUpload)
}