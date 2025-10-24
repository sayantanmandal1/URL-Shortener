package services

import (
	"context"
	"testing"

	"url-shortener-backend/internal/repositories"
	"url-shortener-backend/internal/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyticsService_RecordClick(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	analyticsRepo := repositories.NewAnalyticsRepository(testDB.DB)
	linkRepo := repositories.NewLinkRepository(testDB.DB)
	aiService := testutils.NewMockAIService()
	service := NewAnalyticsService(analyticsRepo, linkRepo, aiService)
	ctx := context.Background()

	// Create a test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	metadata := ClickMetadata{
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0 Test Browser",
		Referrer:  "https://google.com",
	}

	// Test recording click
	err := service.RecordClick(ctx, linkID, metadata)
	require.NoError(t, err)

	// Verify click was recorded
	clicks, err := service.GetClicksByLinkID(ctx, linkID)
	require.NoError(t, err)
	assert.Len(t, clicks, 1)
	assert.Equal(t, linkID, clicks[0].LinkID)
	assert.Equal(t, metadata.IPAddress, clicks[0].IPAddress)
	assert.Equal(t, metadata.UserAgent, clicks[0].UserAgent)
	assert.Equal(t, metadata.Referrer, clicks[0].Referrer)
}

func TestAnalyticsService_RecordClick_NonExistentLink(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	analyticsRepo := repositories.NewAnalyticsRepository(testDB.DB)
	linkRepo := repositories.NewLinkRepository(testDB.DB)
	aiService := testutils.NewMockAIService()
	service := NewAnalyticsService(analyticsRepo, linkRepo, aiService)
	ctx := context.Background()

	metadata := ClickMetadata{
		IPAddress: "192.168.1.1",
		UserAgent: "Mozilla/5.0 Test Browser",
		Referrer:  "",
	}

	// Test recording click for non-existent link
	err := service.RecordClick(ctx, "non-existent-id", metadata)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "link not found")
}

func TestAnalyticsService_GetLinkAnalytics(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	analyticsRepo := repositories.NewAnalyticsRepository(testDB.DB)
	linkRepo := repositories.NewLinkRepository(testDB.DB)
	aiService := testutils.NewMockAIService()
	service := NewAnalyticsService(analyticsRepo, linkRepo, aiService)
	ctx := context.Background()

	// Test analytics for non-existent link
	analytics, err := service.GetLinkAnalytics(ctx, "non-existent-id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "link not found")
	assert.Nil(t, analytics)

	// Create a test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Test analytics for link with no clicks
	analytics, err = service.GetLinkAnalytics(ctx, linkID)
	require.NoError(t, err)
	assert.NotNil(t, analytics)
	assert.Equal(t, 0, analytics.TotalClicks)

	// Record some clicks
	metadata1 := ClickMetadata{
		IPAddress: "192.168.1.1",
		UserAgent: "Desktop Browser",
		Referrer:  "",
	}
	metadata2 := ClickMetadata{
		IPAddress: "192.168.1.2",
		UserAgent: "Mobile Browser",
		Referrer:  "https://google.com",
	}

	require.NoError(t, service.RecordClick(ctx, linkID, metadata1))
	require.NoError(t, service.RecordClick(ctx, linkID, metadata2))

	// Test analytics with clicks
	analytics, err = service.GetLinkAnalytics(ctx, linkID)
	require.NoError(t, err)
	assert.NotNil(t, analytics)
	assert.Equal(t, 2, analytics.TotalClicks)
	assert.NotEmpty(t, analytics.DeviceBreakdown)
	assert.NotEmpty(t, analytics.ReferrerBreakdown)
}

func TestAnalyticsService_GenerateInsights(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	analyticsRepo := repositories.NewAnalyticsRepository(testDB.DB)
	linkRepo := repositories.NewLinkRepository(testDB.DB)
	aiService := testutils.NewMockAIService()
	service := NewAnalyticsService(analyticsRepo, linkRepo, aiService)
	ctx := context.Background()

	// Test insights for non-existent link
	insights, err := service.GenerateInsights(ctx, "non-existent-id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "link not found")
	assert.Empty(t, insights)

	// Create a test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Test insights for link with no clicks
	insights, err = service.GenerateInsights(ctx, linkID)
	require.NoError(t, err)
	assert.Contains(t, insights, "hasn't received any clicks")

	// Record some clicks to get meaningful insights
	metadata := ClickMetadata{
		IPAddress: "192.168.1.1",
		UserAgent: "Desktop Browser",
		Referrer:  "https://google.com",
	}
	require.NoError(t, service.RecordClick(ctx, linkID, metadata))

	// Test insights with clicks (should use AI service)
	insights, err = service.GenerateInsights(ctx, linkID)
	require.NoError(t, err)
	assert.NotEmpty(t, insights)
	assert.Contains(t, insights, "Mock insights") // From our mock AI service
}

func TestAnalyticsService_GetClicksByLinkID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	analyticsRepo := repositories.NewAnalyticsRepository(testDB.DB)
	linkRepo := repositories.NewLinkRepository(testDB.DB)
	aiService := testutils.NewMockAIService()
	service := NewAnalyticsService(analyticsRepo, linkRepo, aiService)
	ctx := context.Background()

	// Test clicks for non-existent link
	clicks, err := service.GetClicksByLinkID(ctx, "non-existent-id")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "link not found")
	assert.Nil(t, clicks)

	// Create a test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Test clicks for link with no clicks
	clicks, err = service.GetClicksByLinkID(ctx, linkID)
	require.NoError(t, err)
	assert.Empty(t, clicks)

	// Record some clicks
	metadata1 := ClickMetadata{
		IPAddress: "192.168.1.1",
		UserAgent: "Browser 1",
		Referrer:  "",
	}
	metadata2 := ClickMetadata{
		IPAddress: "192.168.1.2",
		UserAgent: "Browser 2",
		Referrer:  "https://twitter.com",
	}

	require.NoError(t, service.RecordClick(ctx, linkID, metadata1))
	require.NoError(t, service.RecordClick(ctx, linkID, metadata2))

	// Test clicks with data
	clicks, err = service.GetClicksByLinkID(ctx, linkID)
	require.NoError(t, err)
	assert.Len(t, clicks, 2)
}

func TestAnalyticsService_DeviceTypeExtraction(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	analyticsRepo := repositories.NewAnalyticsRepository(testDB.DB)
	linkRepo := repositories.NewLinkRepository(testDB.DB)
	aiService := testutils.NewMockAIService()
	service := NewAnalyticsService(analyticsRepo, linkRepo, aiService)

	tests := []struct {
		userAgent    string
		expectedType string
	}{
		{
			userAgent:    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X)",
			expectedType: "Mobile",
		},
		{
			userAgent:    "Mozilla/5.0 (iPad; CPU OS 14_7_1 like Mac OS X)",
			expectedType: "Tablet",
		},
		{
			userAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			expectedType: "Desktop",
		},
		{
			userAgent:    "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			expectedType: "Bot",
		},
		{
			userAgent:    "",
			expectedType: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expectedType, func(t *testing.T) {
			deviceType := service.extractDeviceType(tt.userAgent)
			assert.Equal(t, tt.expectedType, deviceType)
		})
	}
}
