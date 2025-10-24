package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"url-shortener-backend/internal/errors"
	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type LinkHandler struct {
	linkService      *services.LinkService
	analyticsService *services.AnalyticsService
}

func NewLinkHandler(linkService *services.LinkService, analyticsService *services.AnalyticsService) *LinkHandler {
	return &LinkHandler{
		linkService:      linkService,
		analyticsService: analyticsService,
	}
}

// CreateLink handles POST /api/links - creates a new shortened link
func (h *LinkHandler) CreateLink(c *fiber.Ctx) error {
	var req models.CreateLinkRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(c, errors.NewValidationError("Invalid JSON format: "+err.Error()))
	}

	// Validate required fields
	if req.OriginalURL == "" {
		return h.handleError(c, errors.NewValidationError("original_url is required and cannot be empty"))
	}

	// Validate custom slug if provided
	if req.CustomSlug != "" {
		if len(req.CustomSlug) < 3 {
			return h.handleError(c, errors.NewValidationError("custom_slug must be at least 3 characters long"))
		}
		if len(req.CustomSlug) > 50 {
			return h.handleError(c, errors.NewValidationError("custom_slug cannot exceed 50 characters"))
		}
	}

	// Create the link
	link, err := h.linkService.CreateLink(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return created link with 201 status
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    link,
		"message": "Link created successfully",
	})
}

// GetLinks handles GET /api/links - retrieves all links
func (h *LinkHandler) GetLinks(c *fiber.Ctx) error {
	links, err := h.linkService.GetAllLinks(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    links,
	})
}

// GetLink handles GET /api/links/:id - retrieves a specific link by ID
func (h *LinkHandler) GetLink(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return h.handleError(c, errors.NewValidationError("Link ID is required"))
	}

	link, err := h.linkService.GetLink(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    link,
	})
}

// RedirectLink handles GET /:slug - redirects to original URL and tracks click
func (h *LinkHandler) RedirectLink(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return h.handleError(c, errors.NewValidationError("Slug is required"))
	}

	// Get the link first
	link, err := h.linkService.GetLinkBySlug(c.Context(), slug)
	if err != nil {
		return h.handleError(c, err)
	}

	// Capture metadata before redirect
	metadata := services.ClickMetadata{
		IPAddress: h.getClientIP(c),
		UserAgent: c.Get("User-Agent"),
		Referrer:  c.Get("Referer"),
	}

	// Record click analytics and increment count asynchronously
	// Use background context to avoid cancellation issues
	go func() {
		ctx := context.Background()

		fmt.Printf("Recording click for link %s (slug: %s) from IP: %s\n", link.ID, slug, metadata.IPAddress)

		// Increment click count
		if err := h.linkService.IncrementClickCount(ctx, link.ID); err != nil {
			fmt.Printf("Error: Failed to increment click count for link %s: %v\n", link.ID, err)
		} else {
			fmt.Printf("Successfully incremented click count for link %s\n", link.ID)
		}

		// Record detailed analytics
		if err := h.analyticsService.RecordClick(ctx, link.ID, metadata); err != nil {
			fmt.Printf("Error: Failed to record click analytics for link %s: %v\n", link.ID, err)
		} else {
			fmt.Printf("Successfully recorded click analytics for link %s\n", link.ID)
		}
	}()

	// Perform redirect immediately
	return c.Redirect(link.OriginalURL, http.StatusFound)
}

// handleError handles application errors and returns appropriate HTTP responses
func (h *LinkHandler) handleError(c *fiber.Ctx, err error) error {
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
func (h *LinkHandler) getStatusCodeFromError(errorCode string) int {
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

// getClientIP extracts the client IP address from the request
func (h *LinkHandler) getClientIP(c *fiber.Ctx) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := c.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to remote address
	return c.IP()
}
