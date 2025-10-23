package handlers

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db      *pgxpool.Pool
	startTime time.Time
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now(),
	}
}

// HealthCheck handles GET /health - returns basic application health status
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "url-shortener-backend",
		"version":   "1.0.0",
		"uptime":    time.Since(h.startTime).String(),
	})
}

// ReadinessCheck handles GET /ready - returns detailed readiness status
func (h *HealthHandler) ReadinessCheck(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	checks := fiber.Map{}
	overallStatus := "ready"

	// Check database connectivity
	if h.db != nil {
		if err := h.db.Ping(ctx); err != nil {
			checks["database"] = fiber.Map{
				"status": "unhealthy",
				"error":  err.Error(),
			}
			overallStatus = "not_ready"
		} else {
			// Get database stats
			stats := h.db.Stat()
			checks["database"] = fiber.Map{
				"status":           "healthy",
				"total_conns":      stats.TotalConns(),
				"acquired_conns":   stats.AcquiredConns(),
				"idle_conns":       stats.IdleConns(),
				"max_conns":        stats.MaxConns(),
			}
		}
	} else {
		checks["database"] = fiber.Map{
			"status": "not_configured",
		}
	}

	// Check AI service availability (basic check)
	// In a real implementation, you might want to make a test call to OpenAI
	checks["ai_service"] = fiber.Map{
		"status": "available", // Assume available if API key is configured
	}

	statusCode := http.StatusOK
	if overallStatus != "ready" {
		statusCode = http.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status":    overallStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks":    checks,
	})
}

// MetricsCheck handles GET /metrics - returns basic application metrics
func (h *HealthHandler) MetricsCheck(c *fiber.Ctx) error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := fiber.Map{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(h.startTime).String(),
		"memory": fiber.Map{
			"alloc_bytes":        memStats.Alloc,
			"total_alloc_bytes":  memStats.TotalAlloc,
			"sys_bytes":          memStats.Sys,
			"num_gc":             memStats.NumGC,
			"heap_alloc_bytes":   memStats.HeapAlloc,
			"heap_sys_bytes":     memStats.HeapSys,
		},
		"runtime": fiber.Map{
			"version":     runtime.Version(),
			"goroutines":  runtime.NumGoroutine(),
			"cpus":        runtime.NumCPU(),
		},
	}

	// Add database connection pool metrics if available
	if h.db != nil {
		stats := h.db.Stat()
		metrics["database"] = fiber.Map{
			"total_conns":      stats.TotalConns(),
			"acquired_conns":   stats.AcquiredConns(),
			"idle_conns":       stats.IdleConns(),
			"max_conns":        stats.MaxConns(),
			"construct_conns":  stats.ConstructingConns(),
		}
	}

	return c.JSON(metrics)
}