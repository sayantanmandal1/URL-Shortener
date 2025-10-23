package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ValidateURL validates if the provided string is a valid URL
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Check if scheme is present and valid
	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must include a scheme (http:// or https://)")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https")
	}

	// Check if host is present
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must include a valid host")
	}

	// Additional validation for common URL issues
	if strings.Contains(parsedURL.Host, " ") {
		return fmt.Errorf("URL host cannot contain spaces")
	}

	return nil
}

// ValidateCustomSlug validates a custom slug provided by the user
func ValidateCustomSlug(slug string) error {
	if slug == "" {
		return nil // Empty slug is allowed (will auto-generate)
	}

	// Check length constraints
	if len(slug) < 3 {
		return fmt.Errorf("custom slug must be at least 3 characters long")
	}

	if len(slug) > 50 {
		return fmt.Errorf("custom slug must be no more than 50 characters long")
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	validSlugPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validSlugPattern.MatchString(slug) {
		return fmt.Errorf("custom slug can only contain letters, numbers, hyphens, and underscores")
	}

	// Check for reserved words/patterns
	reservedWords := []string{
		"api", "admin", "www", "app", "dashboard", "analytics", 
		"health", "status", "docs", "swagger", "openapi",
	}

	lowerSlug := strings.ToLower(slug)
	for _, reserved := range reservedWords {
		if lowerSlug == reserved {
			return fmt.Errorf("slug '%s' is reserved and cannot be used", slug)
		}
	}

	return nil
}