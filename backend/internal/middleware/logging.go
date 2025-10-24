package middleware

import (
	"os"
	"url-shortener-backend/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupLogging configures logging middleware based on environment
func SetupLogging(cfg *config.Config) fiber.Handler {
	var format string

	if cfg.IsProduction() {
		// More detailed logging for production
		format = "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency} - ${userAgent} - ${referer}\n"
	} else {
		// Simpler logging for development
		format = "[${time}] ${status} - ${method} ${path} - ${ip} - ${latency}\n"
	}

	return logger.New(logger.Config{
		Format:     format,
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "UTC",
		Output:     os.Stdout,
	})
}
