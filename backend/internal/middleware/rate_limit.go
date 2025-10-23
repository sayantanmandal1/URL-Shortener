package middleware

import (
	"time"
	"url-shortener-backend/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// createRateLimit creates a rate limiting middleware with the specified configuration
func createRateLimit(max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as the key for rate limiting
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests, please try again later",
				},
			})
		},
	})
}

// DefaultRateLimit provides a default rate limiting configuration based on config
func DefaultRateLimit(cfg *config.Config) fiber.Handler {
	return createRateLimit(cfg.RateLimit.DefaultMax, cfg.RateLimit.DefaultWindow)
}

// StrictRateLimit provides a stricter rate limiting configuration for link creation
func StrictRateLimit(cfg *config.Config) fiber.Handler {
	return createRateLimit(cfg.RateLimit.StrictMax, cfg.RateLimit.StrictWindow)
}

// RedirectRateLimit provides rate limiting for redirect endpoints (more lenient)
func RedirectRateLimit(cfg *config.Config) fiber.Handler {
	return createRateLimit(cfg.RateLimit.RedirectMax, cfg.RateLimit.RedirectWindow)
}