package middleware

import (
	"url-shortener-backend/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupCORS configures CORS middleware based on configuration
func SetupCORS(cfg *config.Config) fiber.Handler {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: cfg.IsProduction(), // Only allow credentials in production
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400, // 24 hours
	}

	return cors.New(corsConfig)
}