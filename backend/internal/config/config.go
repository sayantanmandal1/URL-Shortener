package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	OpenAI    OpenAIConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	URL            string
	MaxConnections int
	MaxIdleTime    time.Duration
	MaxLifetime    time.Duration
}

// OpenAIConfig holds OpenAI API configuration
type OpenAIConfig struct {
	APIKey string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	DefaultMax     int
	DefaultWindow  time.Duration
	StrictMax      int
	StrictWindow   time.Duration
	RedirectMax    int
	RedirectWindow time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8000"),
			Environment:  getEnv("ENVIRONMENT", "development"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			URL:            getEnv("DATABASE_URL", ""),
			MaxConnections: getIntEnv("DB_MAX_CONNECTIONS", 25),
			MaxIdleTime:    getDurationEnv("DB_MAX_IDLE_TIME", 15*time.Minute),
			MaxLifetime:    getDurationEnv("DB_MAX_LIFETIME", time.Hour),
		},
		OpenAI: OpenAIConfig{
			APIKey: getEnv("OPENAI_API_KEY", ""),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
			AllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS"),
			AllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization,X-Requested-With"),
		},
		RateLimit: RateLimitConfig{
			DefaultMax:     getIntEnv("RATE_LIMIT_DEFAULT_MAX", 100),
			DefaultWindow:  getDurationEnv("RATE_LIMIT_DEFAULT_WINDOW", time.Minute),
			StrictMax:      getIntEnv("RATE_LIMIT_STRICT_MAX", 10),
			StrictWindow:   getDurationEnv("RATE_LIMIT_STRICT_WINDOW", time.Minute),
			RedirectMax:    getIntEnv("RATE_LIMIT_REDIRECT_MAX", 1000),
			RedirectWindow: getDurationEnv("RATE_LIMIT_REDIRECT_WINDOW", time.Minute),
		},
	}
}

// IsProduction returns true if the environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// IsDevelopment returns true if the environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
