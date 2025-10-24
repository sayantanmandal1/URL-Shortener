package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener-backend/internal/config"
	"url-shortener-backend/internal/database"
	"url-shortener-backend/internal/handlers"
	"url-shortener-backend/internal/middleware"
	"url-shortener-backend/internal/repositories"
	"url-shortener-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Server represents the HTTP server
type Server struct {
	app    *fiber.App
	config *config.Config
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	// Create Fiber app with configuration
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "HTTP_ERROR",
					"message": err.Error(),
				},
			})
		},
	})

	return &Server{
		app:    app,
		config: cfg,
	}
}

// SetupMiddleware configures all middleware for the application
func (s *Server) SetupMiddleware() {
	// CORS middleware (should be first to handle preflight requests)
	s.app.Use(middleware.SetupCORS(s.config))

	// Recovery middleware
	s.app.Use(recover.New())

	// Logging middleware
	s.app.Use(middleware.SetupLogging(s.config))

	// Default rate limiting for all routes
	s.app.Use(middleware.DefaultRateLimit(s.config))
}

// SetupRoutes configures all routes for the application
func (s *Server) SetupRoutes() error {
	// Initialize database connection
	db, err := database.NewConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Initialize repositories
	linkRepo := repositories.NewLinkRepository(db)
	analyticsRepo := repositories.NewAnalyticsRepository(db)

	// Initialize services
	aiService := services.NewAIService(s.config.OpenAI.APIKey)
	linkService := services.NewLinkService(linkRepo, aiService)
	analyticsService := services.NewAnalyticsService(analyticsRepo, linkRepo, aiService)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db.Pool)
	linkHandler := handlers.NewLinkHandler(linkService, analyticsService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Health check routes (no additional rate limiting)
	s.app.Get("/health", healthHandler.HealthCheck)
	s.app.Get("/ready", healthHandler.ReadinessCheck)
	s.app.Get("/metrics", healthHandler.MetricsCheck)

	// Simple CORS test endpoint
	s.app.Get("/cors-test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "CORS is working",
			"origin":  c.Get("Origin"),
		})
	})

	// API routes with default rate limiting
	api := s.app.Group("/api")

	// Link management routes
	links := api.Group("/links")
	links.Get("/", linkHandler.GetLinks)
	links.Get("/:id", linkHandler.GetLink)

	// Link creation with stricter rate limiting
	links.Post("/", middleware.StrictRateLimit(s.config), linkHandler.CreateLink)

	// Analytics routes
	analytics := api.Group("/analytics")
	analytics.Get("/:id", analyticsHandler.GetAnalytics)
	analytics.Get("/:id/clicks", analyticsHandler.GetClicks)

	// AI insights routes
	insights := api.Group("/insights")
	insights.Get("/:id", analyticsHandler.GetInsights)

	// Test routes (only in development)
	if s.config.IsDevelopment() {
		test := api.Group("/test")
		test.Post("/generate-clicks/:id", analyticsHandler.GenerateTestClicks)
	}

	// Redirect routes with more lenient rate limiting
	s.app.Get("/:slug", middleware.RedirectRateLimit(s.config), linkHandler.RedirectLink)

	return nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Setup middleware
	s.SetupMiddleware()

	// Setup routes
	if err := s.SetupRoutes(); err != nil {
		return fmt.Errorf("failed to setup routes: %w", err)
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", s.config.Server.Port)
		if err := s.app.Listen(":" + s.config.Server.Port); err != nil {
			log.Printf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")

	// Gracefully shutdown the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited")
	return nil
}

// GetApp returns the Fiber app instance (useful for testing)
func (s *Server) GetApp() *fiber.App {
	return s.app
}
