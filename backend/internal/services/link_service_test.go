package services

import (
	"context"
	"testing"

	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/repositories"
	"url-shortener-backend/internal/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkService_CreateLink(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	linkRepo := repositories.NewLinkRepository(testDB.DB)
	aiService := testutils.NewMockAIService()
	service := NewLinkService(linkRepo, aiService)
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
			errorMsg:    "invalid URL",
		},
		{
			name: "custom slug too short",
			request: models.CreateLinkRequest{
				OriginalURL: "https://example.com",
				CustomSlug:  "ab",
			},
			expectError: true,
			errorMsg:    "validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDB.ClearTables(t)

			link, err := service.CreateLink(ctx, tt.request)

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
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	linkRepo := repositories.NewLinkRepository(testDB.DB)
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
	assert.Contains(t, err.Error(), "slug already exists")
	assert.Nil(t, link2)
}

func TestLinkService_GetLink(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	linkRepo := repositories.NewLinkRepository(testDB.DB)
	aiService := testutils.NewMockAIService()
	service := NewLinkService(linkRepo, aiService)
	ctx := context.Background()

	// Test getting non-existent link
	link, err := service.GetLink(ctx, "non-existent-id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "link not found")
	assert.Nil(t, link)

	// Create a test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Test getting existing link
	link, err = service.GetLink(ctx, linkID)
	require.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, linkID, link.ID)
	assert.Equal(t, "https://example.com", link.OriginalURL)
	assert.Equal(t, "test-slug", link.Slug)
}

func TestLinkService_GetLinkBySlug(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	linkRepo := repositories.NewLinkRepository(testDB.DB)
	aiService := testutils.NewMockAIService()
	service := NewLinkService(linkRepo, aiService)
	ctx := context.Background()

	// Test getting non-existent slug
	link, err := service.GetLinkBySlug(ctx, "non-existent-slug")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "link not found")
	assert.Nil(t, link)

	// Create a test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Test getting existing slug
	link, err = service.GetLinkBySlug(ctx, "test-slug")
	require.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, linkID, link.ID)
	assert.Equal(t, "https://example.com", link.OriginalURL)
	assert.Equal(t, "test-slug", link.Slug)
}

func TestLinkService_GetAllLinks(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	linkRepo := repositories.NewLinkRepository(testDB.DB)
	aiService := testutils.NewMockAIService()
	service := NewLinkService(linkRepo, aiService)
	ctx := context.Background()

	// Test empty database
	links, err := service.GetAllLinks(ctx)
	require.NoError(t, err)
	assert.Empty(t, links)

	// Create test links
	linkID1 := uuid.New().String()
	linkID2 := uuid.New().String()
	testDB.CreateTestLink(t, linkID1, "https://example1.com", "slug1")
	testDB.CreateTestLink(t, linkID2, "https://example2.com", "slug2")

	// Test getting all links
	links, err = service.GetAllLinks(ctx)
	require.NoError(t, err)
	assert.Len(t, links, 2)
}

func TestLinkService_IncrementClickCount(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	linkRepo := repositories.NewLinkRepository(testDB.DB)
	aiService := testutils.NewMockAIService()
	service := NewLinkService(linkRepo, aiService)
	ctx := context.Background()

	// Test incrementing non-existent link
	err := service.IncrementClickCount(ctx, "non-existent-id")
	require.Error(t, err)

	// Create a test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Test incrementing existing link
	err = service.IncrementClickCount(ctx, linkID)
	require.NoError(t, err)

	// Verify click count was incremented
	link, err := service.GetLink(ctx, linkID)
	require.NoError(t, err)
	assert.Equal(t, 1, link.ClickCount)
}
