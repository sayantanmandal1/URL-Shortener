package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expectError bool
	}{
		{
			name:        "valid HTTP URL",
			url:         "http://example.com",
			expectError: false,
		},
		{
			name:        "valid HTTPS URL",
			url:         "https://example.com",
			expectError: false,
		},
		{
			name:        "valid URL with path",
			url:         "https://example.com/path/to/resource",
			expectError: false,
		},
		{
			name:        "valid URL with query params",
			url:         "https://example.com/search?q=test&page=1",
			expectError: false,
		},
		{
			name:        "valid URL with fragment",
			url:         "https://example.com/page#section",
			expectError: false,
		},
		{
			name:        "invalid URL - no scheme",
			url:         "example.com",
			expectError: true,
		},
		{
			name:        "invalid URL - empty",
			url:         "",
			expectError: true,
		},
		{
			name:        "invalid URL - malformed",
			url:         "not-a-url",
			expectError: true,
		},
		{
			name:        "invalid URL - unsupported scheme",
			url:         "ftp://example.com",
			expectError: true,
		},
		{
			name:        "invalid URL - localhost",
			url:         "http://localhost:3000",
			expectError: true,
		},
		{
			name:        "invalid URL - private IP",
			url:         "http://192.168.1.1",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCustomSlug(t *testing.T) {
	tests := []struct {
		name        string
		slug        string
		expectError bool
	}{
		{
			name:        "valid slug - letters and numbers",
			slug:        "abc123",
			expectError: false,
		},
		{
			name:        "valid slug - with hyphens",
			slug:        "my-custom-slug",
			expectError: false,
		},
		{
			name:        "valid slug - with underscores",
			slug:        "my_custom_slug",
			expectError: false,
		},
		{
			name:        "valid slug - mixed case",
			slug:        "MyCustomSlug",
			expectError: false,
		},
		{
			name:        "valid slug - minimum length",
			slug:        "abc",
			expectError: false,
		},
		{
			name:        "valid slug - maximum length",
			slug:        "a123456789012345678901234567890123456789012345678",
			expectError: false,
		},
		{
			name:        "invalid slug - too short",
			slug:        "ab",
			expectError: true,
		},
		{
			name:        "invalid slug - too long",
			slug:        "a12345678901234567890123456789012345678901234567890",
			expectError: true,
		},
		{
			name:        "invalid slug - empty",
			slug:        "",
			expectError: true,
		},
		{
			name:        "invalid slug - special characters",
			slug:        "my@slug",
			expectError: true,
		},
		{
			name:        "invalid slug - spaces",
			slug:        "my slug",
			expectError: true,
		},
		{
			name:        "invalid slug - dots",
			slug:        "my.slug",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCustomSlug(tt.slug)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "length 6",
			length: 6,
		},
		{
			name:   "length 8",
			length: 8,
		},
		{
			name:   "length 10",
			length: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug, err := GenerateSlug(tt.length)
			assert.NoError(t, err)
			assert.Len(t, slug, tt.length)

			// Verify slug contains only valid characters
			for _, char := range slug {
				assert.True(t,
					(char >= 'a' && char <= 'z') ||
						(char >= 'A' && char <= 'Z') ||
						(char >= '0' && char <= '9'),
					"Slug contains invalid character: %c", char)
			}
		})
	}
}

func TestGenerateUniqueSlug(t *testing.T) {
	tests := []struct {
		name        string
		customSlug  string
		originalURL string
	}{
		{
			name:        "with custom slug",
			customSlug:  "my-custom",
			originalURL: "https://example.com",
		},
		{
			name:        "without custom slug",
			customSlug:  "",
			originalURL: "https://example.com",
		},
		{
			name:        "with URL-based generation",
			customSlug:  "",
			originalURL: "https://github.com/user/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug, err := GenerateUniqueSlug(tt.customSlug, tt.originalURL)
			assert.NoError(t, err)
			assert.NotEmpty(t, slug)

			if tt.customSlug != "" {
				assert.Equal(t, tt.customSlug, slug)
			} else {
				// Generated slug should be at least 6 characters
				assert.GreaterOrEqual(t, len(slug), 6)
			}
		})
	}
}
