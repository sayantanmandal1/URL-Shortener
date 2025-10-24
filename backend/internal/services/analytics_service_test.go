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

func TestAnalyticsService_RecordClick(t *testing.T) {
	linkRepo := testutils.NewMockLinkRepository()
	analyticsRepo := testutils.NewMockAnalyticsRepository()
	aiService := testutils.NewMockAIService()
	service := NewAnalyticsService(analyticsRepo, linkRepo, aiService)
	ctx := context.Background()

	// Create a test link first
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

	metadata := ClickMetadata{
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0 Test Browser",
		Referrer:  "https://google.com",
	}

	err := service.RecordClick(ctx, testLink.ID, metadata)
	require.NoError(t, err)

	// Verify the click was recorded
	clicks, err := service.GetClicksByLinkID(ctx, testLink.ID)
	require.NoError(t, err)
	assert.Len(t, clicks, 1)
	assert.Equal(t, testLink.ID, clicks[0].LinkID)
	assert.Equal(t, metadata.IPAddress, clicks[0].IPAddress)
	assert.Equal(t, metadata.UserAgent, clicks[0].UserAgent)
	assert.Equal(t, metadata.Referrer, clicks[0].Referrer)
}

func TestAnalyticsService_GetLinkAnalytics(t *testing.T) {
	linkRepo := testutils.NewMockLinkRepository()
	analyticsRepo := testutils.NewMockAnalyticsRepository()
	aiService := testutils.NewMockAIService()
	service := NewAnalyticsService(analyticsRepo, linkRepo, aiService)
	ctx := context.Background()

	// Create a test link first
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

	// Test getting analytics for link with no clicks
	analytics, err := service.GetLinkAnalytics(ctx, testLink.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, analytics.TotalClicks)
	assert.Equal(t, 0, analytics.UniqueVisitors)

	// Record some clicks
	click1 := &models.Click{
		ID:         uuid.New().String(),
		LinkID:     testLink.ID,
		ClickedAt:  time.Now(),
		IPAddress:  "192.168.1.1",
		UserAgent:  "Browser 1",
		Referrer:   "",
		Country:    "US",
		DeviceType: "Desktop",
	}
	click2 := &models.Click{
		ID:         uuid.New().String(),
		LinkID:     testLink.ID,
		ClickedAt:  time.Now(),
		IPAddress:  "192.168.1.2",
		UserAgent:  "Browser 2",
		Referrer:   "https://google.com",
		Country:    "UK",
		DeviceType: "Mobile",
	}

	analyticsRepo.RecordClick(ctx, click1)
	analyticsRepo.RecordClick(ctx, click2)

	// Test getting analytics with clicks
	analytics, err = service.GetLinkAnalytics(ctx, testLink.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, analytics.TotalClicks)
	assert.Equal(t, 2, analytics.UniqueVisitors)
}

func TestAnalyticsService_GenerateInsights(t *testing.T) {
	linkRepo := testutils.NewMockLinkRepository()
	analyticsRepo := testutils.NewMockAnalyticsRepository()
	aiService := testutils.NewMockAIService()
	service := NewAnalyticsService(analyticsRepo, linkRepo, aiService)
	ctx := context.Background()

	// Create a test link first
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

	insights, err := service.GenerateInsights(ctx, testLink.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, insights)
	assert.Contains(t, insights, "Mock AI insights")
}
