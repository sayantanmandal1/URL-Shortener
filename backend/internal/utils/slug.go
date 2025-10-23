package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"strings"
)

const (
	// DefaultSlugLength is the default length for auto-generated slugs
	DefaultSlugLength = 8
	// SlugCharset contains the characters used for slug generation
	SlugCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateSlug creates a random slug of the specified length
func GenerateSlug(length int) (string, error) {
	if length <= 0 {
		length = DefaultSlugLength
	}

	slug := make([]byte, length)
	charsetLen := big.NewInt(int64(len(SlugCharset)))

	for i := range slug {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		slug[i] = SlugCharset[randomIndex.Int64()]
	}

	return string(slug), nil
}

// GenerateSlugFromURL creates a slug based on the URL content (fallback method)
func GenerateSlugFromURL(url string) (string, error) {
	// Extract meaningful parts from URL for slug generation
	// This is a fallback method that creates a more readable slug
	
	// Remove protocol and www
	cleaned := strings.TrimPrefix(url, "https://")
	cleaned = strings.TrimPrefix(cleaned, "http://")
	cleaned = strings.TrimPrefix(cleaned, "www.")
	
	// Take first part of domain and path
	parts := strings.Split(cleaned, "/")
	if len(parts) > 0 {
		domain := parts[0]
		domainParts := strings.Split(domain, ".")
		if len(domainParts) > 0 {
			base := domainParts[0]
			// Limit to first 4 characters and add random suffix
			if len(base) > 4 {
				base = base[:4]
			}
			
			// Generate random suffix
			suffix, err := GenerateSlug(4)
			if err != nil {
				return GenerateSlug(DefaultSlugLength)
			}
			
			return base + suffix, nil
		}
	}
	
	// Fallback to completely random slug
	return GenerateSlug(DefaultSlugLength)
}

// GenerateUniqueSlug generates a slug and provides a mechanism to check uniqueness
// The actual uniqueness check should be done at the service layer with database access
func GenerateUniqueSlug(customSlug string, url string) (string, error) {
	// If custom slug is provided and valid, return it
	if customSlug != "" {
		if err := ValidateCustomSlug(customSlug); err != nil {
			return "", err
		}
		return customSlug, nil
	}

	// Try to generate a meaningful slug from URL first
	slug, err := GenerateSlugFromURL(url)
	if err != nil {
		// Fallback to random slug
		return GenerateSlug(DefaultSlugLength)
	}

	return slug, nil
}

// GenerateBase64Slug creates a URL-safe base64 encoded slug
func GenerateBase64Slug(length int) (string, error) {
	if length <= 0 {
		length = DefaultSlugLength
	}

	// Generate random bytes
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode to base64 and make URL-safe
	slug := base64.URLEncoding.EncodeToString(bytes)
	
	// Remove padding and limit to desired length
	slug = strings.TrimRight(slug, "=")
	if len(slug) > length {
		slug = slug[:length]
	}

	return slug, nil
}