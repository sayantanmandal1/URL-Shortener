package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"url-shortener-backend/internal/errors"
	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/repositories"
	"url-shortener-backend/internal/utils"

	"github.com/google/uuid"
)

type LinkService struct {
	linkRepo  *repositories.LinkRepository
	aiService AIService
}

func NewLinkService(linkRepo *repositories.LinkRepository, aiService AIService) *LinkService {
	return &LinkService{
		linkRepo:  linkRepo,
		aiService: aiService,
	}
}

// CreateLink creates a new shortened link with validation and AI-generated metadata
func (s *LinkService) CreateLink(ctx context.Context, req models.CreateLinkRequest) (*models.Link, error) {
	// Validate the original URL
	if err := utils.ValidateURL(req.OriginalURL); err != nil {
		return nil, errors.NewInvalidURLError(err.Error())
	}

	// Generate or validate the slug
	slug, err := s.generateUniqueSlug(ctx, req.CustomSlug, req.OriginalURL)
	if err != nil {
		return nil, err
	}

	// Generate AI metadata (title and description)
	aiMetadata, err := s.aiService.GenerateTitleAndDescription(ctx, req.OriginalURL)
	if err != nil {
		// Log the error but don't fail the request - use fallback values
		fmt.Printf("Warning: Failed to generate AI metadata: %v\n", err)
	}

	// Create the link entity
	link := &models.Link{
		ID:          uuid.New().String(),
		OriginalURL: req.OriginalURL,
		Slug:        slug,
		Title:       getTitle(aiMetadata, req.OriginalURL),
		Description: getDescription(aiMetadata, req.OriginalURL),
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	if err := s.linkRepo.Create(ctx, link); err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, errors.NewSlugExistsError(slug)
		}
		return nil, errors.NewDatabaseError(err.Error())
	}

	return link, nil
}

// GetLink retrieves a link by its ID
func (s *LinkService) GetLink(ctx context.Context, id string) (*models.Link, error) {
	if id == "" {
		return nil, errors.NewValidationError("Link ID cannot be empty")
	}

	link, err := s.linkRepo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.NewLinkNotFoundError(id)
		}
		return nil, errors.NewDatabaseError(err.Error())
	}

	return link, nil
}

// GetLinkBySlug retrieves a link by its slug
func (s *LinkService) GetLinkBySlug(ctx context.Context, slug string) (*models.Link, error) {
	if slug == "" {
		return nil, errors.NewValidationError("Slug cannot be empty")
	}

	link, err := s.linkRepo.GetBySlug(ctx, slug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.NewLinkNotFoundError(slug)
		}
		return nil, errors.NewDatabaseError(err.Error())
	}

	return link, nil
}

// GetAllLinks retrieves all links ordered by creation date
func (s *LinkService) GetAllLinks(ctx context.Context) ([]*models.Link, error) {
	links, err := s.linkRepo.GetAll(ctx)
	if err != nil {
		return nil, errors.NewDatabaseError(err.Error())
	}

	return links, nil
}

// IncrementClickCount increments the click count for a link
func (s *LinkService) IncrementClickCount(ctx context.Context, linkID string) error {
	return s.linkRepo.IncrementClicks(ctx, linkID)
}

// RedirectAndTrack handles URL redirection (deprecated - use direct methods instead)
func (s *LinkService) RedirectAndTrack(ctx context.Context, slug string) (string, error) {
	// Get the link by slug
	link, err := s.GetLinkBySlug(ctx, slug)
	if err != nil {
		return "", err
	}

	return link.OriginalURL, nil
}

// generateUniqueSlug generates a unique slug, checking for conflicts
func (s *LinkService) generateUniqueSlug(ctx context.Context, customSlug, originalURL string) (string, error) {
	var slug string
	var err error

	// If custom slug is provided, validate and use it
	if customSlug != "" {
		if err := utils.ValidateCustomSlug(customSlug); err != nil {
			return "", errors.NewValidationError(err.Error())
		}
		slug = customSlug
	} else {
		// Generate a slug from the URL or random
		slug, err = utils.GenerateUniqueSlug("", originalURL)
		if err != nil {
			return "", errors.NewAppError(errors.ErrInternalError, "Failed to generate slug", err.Error())
		}
	}

	// Check if slug already exists
	maxAttempts := 5
	for attempt := 0; attempt < maxAttempts; attempt++ {
		_, err := s.linkRepo.GetBySlug(ctx, slug)
		if err != nil {
			// If link not found, slug is available
			if strings.Contains(err.Error(), "not found") {
				return slug, nil
			}
			// Other database error
			return "", errors.NewDatabaseError(err.Error())
		}

		// Slug exists, need to generate a new one
		if customSlug != "" {
			// Custom slug conflicts, return error
			return "", errors.NewSlugExistsError(customSlug)
		}

		// Generate a new random slug for auto-generated slugs
		slug, err = utils.GenerateSlug(8 + attempt) // Increase length with each attempt
		if err != nil {
			return "", errors.NewAppError(errors.ErrInternalError, "Failed to generate slug", err.Error())
		}
	}

	return "", errors.NewAppError(errors.ErrInternalError, "Failed to generate unique slug after multiple attempts", "")
}

// Helper functions for AI metadata fallbacks
func getTitle(aiMetadata *models.AIMetadata, originalURL string) string {
	if aiMetadata != nil && aiMetadata.Title != "" {
		return aiMetadata.Title
	}
	// Fallback: extract domain from URL
	return extractDomainFromURL(originalURL)
}

func getDescription(aiMetadata *models.AIMetadata, originalURL string) string {
	if aiMetadata != nil && aiMetadata.Description != "" {
		return aiMetadata.Description
	}
	// Fallback: generic description
	return fmt.Sprintf("Shortened link to %s", extractDomainFromURL(originalURL))
}

func extractDomainFromURL(rawURL string) string {
	// Simple domain extraction for fallback
	url := strings.TrimPrefix(rawURL, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "www.")

	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return "Unknown"
}
