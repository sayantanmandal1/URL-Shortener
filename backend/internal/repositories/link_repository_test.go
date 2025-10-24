package repositories

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

func TestLinkRepository_Create(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	repo := NewLinkRepository(testDB.DB)
	ctx := context.Background()

	link := &models.Link{
		ID:          uuid.New().String(),
		OriginalURL: "https://example.com",
		Slug:        "test-slug",
		Title:       "Test Title",
		Description: "Test Description",
		ClickCount:  0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.Create(ctx, link)
	require.NoError(t, err)

	// Verify the link was created
	retrieved, err := repo.GetByID(ctx, link.ID)
	require.NoError(t, err)
	assert.Equal(t, link.ID, retrieved.ID)
	assert.Equal(t, link.OriginalURL, retrieved.OriginalURL)
	assert.Equal(t, link.Slug, retrieved.Slug)
	assert.Equal(t, link.Title, retrieved.Title)
	assert.Equal(t, link.Description, retrieved.Description)
	assert.Equal(t, link.ClickCount, retrieved.ClickCount)
}

func TestLinkRepository_GetByID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	repo := NewLinkRepository(testDB.DB)
	ctx := context.Background()

	// Test getting non-existent link
	_, err := repo.GetByID(ctx, "non-existent-id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "link not found")

	// Create a test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Test getting existing link
	link, err := repo.GetByID(ctx, linkID)
	require.NoError(t, err)
	assert.Equal(t, linkID, link.ID)
	assert.Equal(t, "https://example.com", link.OriginalURL)
	assert.Equal(t, "test-slug", link.Slug)
}

func TestLinkRepository_GetBySlug(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	repo := NewLinkRepository(testDB.DB)
	ctx := context.Background()

	// Test getting non-existent slug
	_, err := repo.GetBySlug(ctx, "non-existent-slug")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "link not found")

	// Create a test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Test getting existing slug
	link, err := repo.GetBySlug(ctx, "test-slug")
	require.NoError(t, err)
	assert.Equal(t, linkID, link.ID)
	assert.Equal(t, "https://example.com", link.OriginalURL)
	assert.Equal(t, "test-slug", link.Slug)
}

func TestLinkRepository_GetAll(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	repo := NewLinkRepository(testDB.DB)
	ctx := context.Background()

	// Test empty database
	links, err := repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Empty(t, links)

	// Create test links
	linkID1 := uuid.New().String()
	linkID2 := uuid.New().String()
	testDB.CreateTestLink(t, linkID1, "https://example1.com", "slug1")
	testDB.CreateTestLink(t, linkID2, "https://example2.com", "slug2")

	// Test getting all links
	links, err = repo.GetAll(ctx)
	require.NoError(t, err)
	assert.Len(t, links, 2)

	// Verify links are ordered by created_at DESC
	assert.Equal(t, linkID2, links[0].ID) // More recent
	assert.Equal(t, linkID1, links[1].ID) // Older
}

func TestLinkRepository_IncrementClicks(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	repo := NewLinkRepository(testDB.DB)
	ctx := context.Background()

	// Test incrementing non-existent link
	err := repo.IncrementClicks(ctx, "non-existent-id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "link not found")

	// Create a test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Verify initial click count
	link, err := repo.GetByID(ctx, linkID)
	require.NoError(t, err)
	assert.Equal(t, 0, link.ClickCount)

	// Increment clicks
	err = repo.IncrementClicks(ctx, linkID)
	require.NoError(t, err)

	// Verify click count increased
	link, err = repo.GetByID(ctx, linkID)
	require.NoError(t, err)
	assert.Equal(t, 1, link.ClickCount)

	// Increment again
	err = repo.IncrementClicks(ctx, linkID)
	require.NoError(t, err)

	// Verify click count increased again
	link, err = repo.GetByID(ctx, linkID)
	require.NoError(t, err)
	assert.Equal(t, 2, link.ClickCount)
}
