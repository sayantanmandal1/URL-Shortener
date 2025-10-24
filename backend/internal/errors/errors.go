package errors

import "fmt"

// AppError represents a structured application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Error codes
const (
	ErrInvalidURL      = "INVALID_URL"
	ErrSlugExists      = "SLUG_EXISTS"
	ErrLinkNotFound    = "LINK_NOT_FOUND"
	ErrDatabaseError   = "DATABASE_ERROR"
	ErrAIServiceError  = "AI_SERVICE_ERROR"
	ErrValidationError = "VALIDATION_ERROR"
	ErrInternalError   = "INTERNAL_ERROR"
	ErrInvalidInput    = "INVALID_INPUT"
)

// NewAppError creates a new application error
func NewAppError(code, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    ErrValidationError,
		Message: message,
	}
}

// NewInvalidURLError creates an invalid URL error
func NewInvalidURLError(details string) *AppError {
	return &AppError{
		Code:    ErrInvalidURL,
		Message: "The provided URL is invalid",
		Details: details,
	}
}

// NewSlugExistsError creates a slug exists error
func NewSlugExistsError(slug string) *AppError {
	return &AppError{
		Code:    ErrSlugExists,
		Message: "The requested slug is already in use",
		Details: fmt.Sprintf("Slug '%s' already exists", slug),
	}
}

// NewLinkNotFoundError creates a link not found error
func NewLinkNotFoundError(identifier string) *AppError {
	return &AppError{
		Code:    ErrLinkNotFound,
		Message: "The requested link was not found",
		Details: fmt.Sprintf("Link with identifier '%s' not found", identifier),
	}
}

// NewDatabaseError creates a database error
func NewDatabaseError(details string) *AppError {
	return &AppError{
		Code:    ErrDatabaseError,
		Message: "A database error occurred",
		Details: details,
	}
}

// NewAIServiceError creates an AI service error
func NewAIServiceError(details string) *AppError {
	return &AppError{
		Code:    ErrAIServiceError,
		Message: "AI service error occurred",
		Details: details,
	}
}
