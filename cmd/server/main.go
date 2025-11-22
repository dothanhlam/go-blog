package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go-blog/internal/api"
	"go-blog/internal/config"
	i18nmiddleware "go-blog/internal/middleware"
	"go-blog/internal/service"
	"go-blog/internal/store/postgres"
	"go-blog/internal/storage"
	"go-blog/internal/web"
	"golang.org/x/text/language"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Initialize database connection
	db, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	// Initialize stores
	userStore := postgres.NewUserStore(db)
	postStore := postgres.NewPostStore(db)

	// Initialize file storage
	fileStorage, err := storage.New(cfg)
	if err != nil {
		log.Fatalf("could not initialize file storage: %v", err)
	}

	// Initialize services
	userService := service.NewUserService(userStore)
	postService := service.NewPostService(postStore, fileStorage)

	// Initialize Echo
	e := echo.New()
	e.Validator = api.NewValidator()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Serve static files
	e.Static("/static", "internal/web/static")
	e.Static("/assets", "internal/web/assets")

	// --- Web Routes and Middleware ---
	// The i18n middleware should come before the auth middleware
	e.Use(i18nmiddleware.I18n(language.English))
	e.Use(i18nmiddleware.WebAuth(userService, cfg))

	webHandler := api.NewWebHandler(cfg, postService, userService)
	e.HTTPErrorHandler = webHandler.CustomHTTPErrorHandler

	e.Renderer = web.NewTemplateRenderer()
	e.GET("/posts/:id", webHandler.RenderPostPage)
	e.GET("/", webHandler.RenderIndexPage)
	//e.GET("/login", webHandler.RenderLoginPage)
	//e.POST("/login", webHandler.HandleLogin)
	//e.GET("/logout", webHandler.HandleLogout)

	// Register routes
	api.RegisterRoutes(e, userService, postService, cfg)

	// Start server
	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
}