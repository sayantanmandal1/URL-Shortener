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

func TestAnalyticsRepository_RecordClick(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	repo := NewAnalyticsRepository(testDB.DB)
	ctx := context.Background()

	// Create a test link first
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	click := &models.Click{
		ID:         uuid.New().String(),
		LinkID:     linkID,
		ClickedAt:  time.Now(),
		IPAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0 Test Browser",
		Referrer:   "https://google.com",
		Country:    "US",
		DeviceType: "Desktop",
	}

	err := repo.RecordClick(ctx, click)
	require.NoError(t, err)

	// Verify the click was recorded
	clicks, err := repo.GetClicksByLinkID(ctx, linkID)
	require.NoError(t, err)
	assert.Len(t, clicks, 1)
	assert.Equal(t, click.ID, clicks[0].ID)
	assert.Equal(t, click.LinkID, clicks[0].LinkID)
	assert.Equal(t, click.IPAddress, clicks[0].IPAddress)
	assert.Equal(t, click.UserAgent, clicks[0].UserAgent)
	assert.Equal(t, click.Referrer, clicks[0].Referrer)
	assert.Equal(t, click.Country, clicks[0].Country)
	assert.Equal(t, click.DeviceType, clicks[0].DeviceType)
}

func TestAnalyticsRepository_GetClicksByLinkID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	repo := NewAnalyticsRepository(testDB.DB)
	ctx := context.Background()

	// Create test links
	linkID1 := uuid.New().String()
	linkID2 := uuid.New().String()
	testDB.CreateTestLink(t, linkID1, "https://example1.com", "slug1")
	testDB.CreateTestLink(t, linkID2, "https://example2.com", "slug2")

	// Test getting clicks for link with no clicks
	clicks, err := repo.GetClicksByLinkID(ctx, linkID1)
	require.NoError(t, err)
	assert.Empty(t, clicks)

	// Create clicks for linkID1
	click1 := &models.Click{
		ID:         uuid.New().String(),
		LinkID:     linkID1,
		ClickedAt:  time.Now().Add(-2 * time.Hour),
		IPAddress:  "192.168.1.1",
		UserAgent:  "Browser 1",
		Referrer:   "",
		Country:    "US",
		DeviceType: "Desktop",
	}

	click2 := &models.Click{
		ID:         uuid.New().String(),
		LinkID:     linkID1,
		ClickedAt:  time.Now().Add(-1 * time.Hour),
		IPAddress:  "192.168.1.2",
		UserAgent:  "Browser 2",
		Referrer:   "https://twitter.com",
		Country:    "UK",
		DeviceType: "Mobile",
	}

	// Create click for linkID2
	click3 := &models.Click{
		ID:         uuid.New().String(),
		LinkID:     linkID2,
		ClickedAt:  time.Now(),
		IPAddress:  "192.168.1.3",
		UserAgent:  "Browser 3",
		Referrer:   "",
		Country:    "CA",
		DeviceType: "Tablet",
	}

	require.NoError(t, repo.RecordClick(ctx, click1))
	require.NoError(t, repo.RecordClick(ctx, click2))
	require.NoError(t, repo.RecordClick(ctx, click3))

	// Test getting clicks for linkID1
	clicks, err = repo.GetClicksByLinkID(ctx, linkID1)
	require.NoError(t, err)
	assert.Len(t, clicks, 2)

	// Verify clicks are ordered by clicked_at DESC (most recent first)
	assert.Equal(t, click2.ID, clicks[0].ID)
	assert.Equal(t, click1.ID, clicks[1].ID)

	// Test getting clicks for linkID2
	clicks, err = repo.GetClicksByLinkID(ctx, linkID2)
	require.NoError(t, err)
	assert.Len(t, clicks, 1)
	assert.Equal(t, click3.ID, clicks[0].ID)
}

func TestAnalyticsRepository_GetAnalyticsSummary(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	defer testDB.Cleanup(t)

	repo := NewAnalyticsRepository(testDB.DB)
	ctx := context.Background()

	// Create a test link
	linkID := uuid.New().String()
	testDB.CreateTestLink(t, linkID, "https://example.com", "test-slug")

	// Test analytics for link with no clicks
	analytics, err := repo.GetAnalyticsSummary(ctx, linkID)
	require.NoError(t, err)
	assert.Equal(t, 0, analytics.TotalClicks)
	assert.Empty(t, analytics.DailyClicks)
	assert.Empty(t, analytics.CountryBreakdown)
	assert.Empty(t, analytics.DeviceBreakdown)
	assert.Empty(t, analytics.ReferrerBreakdown)

	// Create test clicks with various data
	clicks := []*models.Click{
		{
			ID:         uuid.New().String(),
			LinkID:     linkID,
			ClickedAt:  time.Now().Add(-1 * time.Hour),
			IPAddress:  "192.168.1.1",
			UserAgent:  "Desktop Browser",
			Referrer:   "",
			Country:    "US",
			DeviceType: "Desktop",
		},
		{
			ID:         uuid.New().String(),
			LinkID:     linkID,
			ClickedAt:  time.Now().Add(-2 * time.Hour),
			IPAddress:  "192.168.1.2",
			UserAgent:  "Mobile Browser",
			Referrer:   "https://google.com",
			Country:    "US",
			DeviceType: "Mobile",
		},
		{
			ID:         uuid.New().String(),
			LinkID:     linkID,
			ClickedAt:  time.Now().Add(-3 * time.Hour),
			IPAddress:  "192.168.1.3",
			UserAgent:  "Desktop Browser",
			Referrer:   "https://twitter.com",
			Country:    "UK",
			DeviceType: "Desktop",
		},
	}

	for _, click := range clicks {
		require.NoError(t, repo.RecordClick(ctx, click))
	}

	// Test analytics with clicks
	analytics, err = repo.GetAnalyticsSummary(ctx, linkID)
	require.NoError(t, err)

	// Verify total clicks
	assert.Equal(t, 3, analytics.TotalClicks)

	// Verify daily clicks (should have data for today)
	assert.NotEmpty(t, analytics.DailyClicks)

	// Verify country breakdown
	assert.Len(t, analytics.CountryBreakdown, 2)
	usClicks := 0
	ukClicks := 0
	for _, country := range analytics.CountryBreakdown {
		if country.Country == "US" {
			usClicks = country.Count
		} else if country.Country == "UK" {
			ukClicks = country.Count
		}
	}
	assert.Equal(t, 2, usClicks)
	assert.Equal(t, 1, ukClicks)

	// Verify device breakdown
	assert.Len(t, analytics.DeviceBreakdown, 2)
	desktopClicks := 0
	mobileClicks := 0
	for _, device := range analytics.DeviceBreakdown {
		if device.DeviceType == "Desktop" {
			desktopClicks = device.Count
		} else if device.DeviceType == "Mobile" {
			mobileClicks = device.Count
		}
	}
	assert.Equal(t, 2, desktopClicks)
	assert.Equal(t, 1, mobileClicks)

	// Verify referrer breakdown
	assert.Len(t, analytics.ReferrerBreakdown, 3)
	directClicks := 0
	googleClicks := 0
	twitterClicks := 0
	for _, referrer := range analytics.ReferrerBreakdown {
		switch referrer.Referrer {
		case "Direct":
			directClicks = referrer.Count
		case "https://google.com":
			googleClicks = referrer.Count
		case "https://twitter.com":
			twitterClicks = referrer.Count
		}
	}
	assert.Equal(t, 1, directClicks)
	assert.Equal(t, 1, googleClicks)
	assert.Equal(t, 1, twitterClicks)
}
