package services

import (
	"context"
	"testing"
	"time"

	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkService_CreateLink(t *testing.T) {
	aiService := testutils.NewMockAIService()
	ctx := context.Background()

	tests := []struct {
		name        string
		request     models.CreateLinkRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid URL without custom slug",
			request: models.CreateLinkRequest{
				OriginalURL: "https://example.com",
			},
			expectError: false,
		},
		{
			name: "valid URL with custom slug",
			request: models.CreateLinkRequest{
				OriginalURL: "https://example.com",
				CustomSlug:  "my-custom-slug",
			},
			expectError: false,
		},
		{
			name: "invalid URL",
			request: models.CreateLinkRequest{
				OriginalURL: "not-a-valid-url",
			},
			expectError: true,
			errorMsg:    "INVALID_URL",
		},
		{
			name: "custom slug too short",
			request: models.CreateLinkRequest{
				OriginalURL: "https://example.com",
				CustomSlug:  "ab",
			},
			expectError: true,
			errorMsg:    "VALIDATION_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the mock repository for each test
			testLinkRepo := testutils.NewMockLinkRepository()
			testService := NewLinkService(testLinkRepo, aiService)

			link, err := testService.CreateLink(ctx, tt.request)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, link)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, link)
				assert.NotEmpty(t, link.ID)
				assert.Equal(t, tt.request.OriginalURL, link.OriginalURL)
				assert.NotEmpty(t, link.Slug)
				assert.NotEmpty(t, link.Title)
				assert.NotEmpty(t, link.Description)
				assert.Equal(t, 0, link.ClickCount)

				if tt.request.CustomSlug != "" {
					assert.Equal(t, tt.request.CustomSlug, link.Slug)
				}
			}
		})
	}
}

func TestLinkService_CreateLink_DuplicateSlug(t *testing.T) {
	linkRepo := testutils.NewMockLinkRepository()
	aiService := testutils.NewMockAIService()
	service := NewLinkService(linkRepo, aiService)
	ctx := context.Background()

	// Create first link with custom slug
	req1 := models.CreateLinkRequest{
		OriginalURL: "https://example1.com",
		CustomSlug:  "duplicate-slug",
	}

	link1, err := service.CreateLink(ctx, req1)
	require.NoError(t, err)
	assert.Equal(t, "duplicate-slug", link1.Slug)

	// Try to create second link with same custom slug
	req2 := models.CreateLinkRequest{
		OriginalURL: "https://example2.com",
		CustomSlug:  "duplicate-slug",
	}

	link2, err := service.CreateLink(ctx, req2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "SLUG_EXISTS")
	assert.Nil(t, link2)
}

func TestLinkService_GetLink(t *testing.T) {
	linkRepo := testutils.NewMockLinkRepository()
	aiService := testutils.NewMockAIService()
	service := NewLinkService(linkRepo, aiService)
	ctx := context.Background()

	// Test getting non-existent link
	link, err := service.GetLink(ctx, "non-existent-id")
	require.Error(t, err)
	assert.Nil(t, link)

	// Create a test link
	testLink := &models.Link{
		ID:          uuid.New().String(),
		OriginalURL: "https://example.com",
		Slug:        "test-slug",
		Title:       "Test Title",
		Description: "Test Description",
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	linkRepo.Create(ctx, testLink)

	// Test getting existing link
	link, err = service.GetLink(ctx, testLink.ID)
	require.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, testLink.ID, link.ID)
	assert.Equal(t, "https://example.com", link.OriginalURL)
	assert.Equal(t, "test-slug", link.Slug)
}

func TestLinkService_GetLinkBySlug(t *testing.T) {
	linkRepo := testutils.NewMockLinkRepository()
	aiService := testutils.NewMockAIService()
	service := NewLinkService(linkRepo, aiService)
	ctx := context.Background()

	// Test getting non-existent slug
	link, err := service.GetLinkBySlug(ctx, "non-existent-slug")
	require.Error(t, err)
	assert.Nil(t, link)

	// Create a test link
	testLink := &models.Link{
		ID:          uuid.New().String(),
		OriginalURL: "https://example.com",
		Slug:        "test-slug",
		Title:       "Test Title",
		Description: "Test Description",
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	linkRepo.Create(ctx, testLink)

	// Test getting existing slug
	link, err = service.GetLinkBySlug(ctx, "test-slug")
	require.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, testLink.ID, link.ID)
	assert.Equal(t, "https://example.com", link.OriginalURL)
	assert.Equal(t, "test-slug", link.Slug)
}

func TestLinkService_GetAllLinks(t *testing.T) {
	linkRepo := testutils.NewMockLinkRepository()
	aiService := testutils.NewMockAIService()
	service := NewLinkService(linkRepo, aiService)
	ctx := context.Background()

	// Test empty repository
	links, err := service.GetAllLinks(ctx)
	require.NoError(t, err)
	assert.Empty(t, links)

	// Create test links
	testLink1 := &models.Link{
		ID:          uuid.New().String(),
		OriginalURL: "https://example1.com",
		Slug:        "slug1",
		Title:       "Test Title 1",
		Description: "Test Description 1",
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	testLink2 := &models.Link{
		ID:          uuid.New().String(),
		OriginalURL: "https://example2.com",
		Slug:        "slug2",
		Title:       "Test Title 2",
		Description: "Test Description 2",
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	linkRepo.Create(ctx, testLink1)
	linkRepo.Create(ctx, testLink2)

	// Test getting all links
	links, err = service.GetAllLinks(ctx)
	require.NoError(t, err)
	assert.Len(t, links, 2)
}

func TestLinkService_IncrementClickCount(t *testing.T) {
	linkRepo := testutils.NewMockLinkRepository()
	aiService := testutils.NewMockAIService()
	service := NewLinkService(linkRepo, aiService)
	ctx := context.Background()

	// Test incrementing non-existent link
	err := service.IncrementClickCount(ctx, "non-existent-id")
	require.Error(t, err)

	// Create a test link
	testLink := &models.Link{
		ID:          uuid.New().String(),
		OriginalURL: "https://example.com",
		Slug:        "test-slug",
		Title:       "Test Title",
		Description: "Test Description",
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	linkRepo.Create(ctx, testLink)

	// Test incrementing existing link
	err = service.IncrementClickCount(ctx, testLink.ID)
	require.NoError(t, err)

	// Verify click count was incremented
	link, err := service.GetLink(ctx, testLink.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, link.ClickCount)
}
