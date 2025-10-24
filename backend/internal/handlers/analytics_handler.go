package handlers

import (
	"net/http"

	"url-shortener-backend/internal/errors"
	"url-shortener-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetAnalytics handles GET /api/analytics/:id - retrieves analytics data for a specific link
func (h *AnalyticsHandler) GetAnalytics(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.NewValidationError("Link ID is required"))
	}

	analytics, err := h.analyticsService.GetLinkAnalytics(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    analytics,
	})
}

// GetInsights handles GET /api/insights/:id - retrieves AI-generated insights for a specific link
func (h *AnalyticsHandler) GetInsights(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.NewValidationError("Link ID is required"))
	}

	insights, err := h.analyticsService.GenerateInsights(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"insights": insights,
		},
	})
}

// GetClicks handles GET /api/analytics/:id/clicks - retrieves all clicks for a specific link
func (h *AnalyticsHandler) GetClicks(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.NewValidationError("Link ID is required"))
	}

	clicks, err := h.analyticsService.GetClicksByLinkID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    clicks,
	})
}

// GenerateTestClicks generates sample click data for testing (development only)
func (h *AnalyticsHandler) GenerateTestClicks(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.NewValidationError("Link ID is required"))
	}

	err := h.analyticsService.GenerateTestClicks(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Test clicks generated successfully",
	})
}

// handleError handles application errors and returns appropriate HTTP responses
func (h *AnalyticsHandler) handleError(c *fiber.Ctx, err error) error {
	// Check if it's an application error
	if appErr, ok := err.(*errors.AppError); ok {
		statusCode := h.getStatusCodeFromError(appErr.Code)
		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    appErr.Code,
				"message": appErr.Message,
				"details": appErr.Details,
			},
		})
	}

	// Handle unknown errors
	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    "INTERNAL_ERROR",
			"message": "An internal error occurred",
			"details": err.Error(),
		},
	})
}

// getStatusCodeFromError maps error codes to HTTP status codes
func (h *AnalyticsHandler) getStatusCodeFromError(errorCode string) int {
	switch errorCode {
	case errors.ErrInvalidURL, errors.ErrValidationError, errors.ErrInvalidInput:
		return http.StatusBadRequest
	case errors.ErrLinkNotFound:
		return http.StatusNotFound
	case errors.ErrSlugExists:
		return http.StatusConflict
	case errors.ErrAIServiceError:
		return http.StatusBadGateway
	case errors.ErrDatabaseError, errors.ErrInternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
