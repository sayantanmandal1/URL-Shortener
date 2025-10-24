package middleware

import (
	"url-shortener-backend/internal/config"

	"github.com/gofiber/fiber/v2"
)

// SetupCORS configures CORS middleware based on configuration
func SetupCORS(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set CORS headers for all requests
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Requested-With,Access-Control-Request-Method,Access-Control-Request-Headers")
		c.Set("Access-Control-Expose-Headers", "Content-Length")
		c.Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
