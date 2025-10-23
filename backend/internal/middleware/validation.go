package middleware

import (
	"net/http"
	"strings"

	"url-shortener-backend/internal/errors"

	"github.com/gofiber/fiber/v2"
)

// ValidateCreateLinkRequest validates the request body for creating a link
func ValidateCreateLinkRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check content type
		contentType := c.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    errors.ErrInvalidInput,
					"message": "Content-Type must be application/json",
				},
			})
		}

		// Check if body is empty
		if len(c.Body()) == 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    errors.ErrInvalidInput,
					"message": "Request body cannot be empty",
				},
			})
		}

		return c.Next()
	}
}

// ValidateIDParam validates that the ID parameter is present and not empty
func ValidateIDParam() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    errors.ErrValidationError,
					"message": "ID parameter is required",
				},
			})
		}

		// Basic UUID format validation (optional but good practice)
		if len(id) < 8 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    errors.ErrValidationError,
					"message": "Invalid ID format",
				},
			})
		}

		return c.Next()
	}
}

// ValidateSlugParam validates that the slug parameter is present and not empty
func ValidateSlugParam() fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		if slug == "" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    errors.ErrValidationError,
					"message": "Slug parameter is required",
				},
			})
		}

		// Basic slug validation
		if len(slug) > 50 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    errors.ErrValidationError,
					"message": "Slug is too long (maximum 50 characters)",
				},
			})
		}

		return c.Next()
	}
}

// ErrorHandler is a global error handler for the application
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Default to 500 server error
		code := http.StatusInternalServerError
		message := "Internal Server Error"

		// Check if it's a Fiber error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message
		}

		// Check if it's an application error
		if appErr, ok := err.(*errors.AppError); ok {
			code = getStatusCodeFromError(appErr.Code)
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    appErr.Code,
					"message": appErr.Message,
					"details": appErr.Details,
				},
			})
		}

		// Return generic error
		return c.Status(code).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": message,
			},
		})
	}
}

// getStatusCodeFromError maps error codes to HTTP status codes
func getStatusCodeFromError(errorCode string) int {
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