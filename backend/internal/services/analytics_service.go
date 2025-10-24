package services

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"url-shortener-backend/internal/errors"
	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/repositories"

	"github.com/google/uuid"
)

type AnalyticsService struct {
	analyticsRepo *repositories.AnalyticsRepository
	linkRepo      *repositories.LinkRepository
	aiService     AIService
}

func NewAnalyticsService(analyticsRepo *repositories.AnalyticsRepository, linkRepo *repositories.LinkRepository, aiService AIService) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
		linkRepo:      linkRepo,
		aiService:     aiService,
	}
}

// ClickMetadata represents metadata captured during a click event
type ClickMetadata struct {
	IPAddress string
	UserAgent string
	Referrer  string
}

// RecordClick records a click event with metadata capture
func (s *AnalyticsService) RecordClick(ctx context.Context, linkID string, metadata ClickMetadata) error {
	if linkID == "" {
		return errors.NewValidationError("Link ID cannot be empty")
	}

	// Verify the link exists
	_, err := s.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.NewLinkNotFoundError(linkID)
		}
		return errors.NewDatabaseError(err.Error())
	}

	// Create click record
	click := &models.Click{
		ID:         uuid.New().String(),
		LinkID:     linkID,
		ClickedAt:  time.Now(),
		IPAddress:  metadata.IPAddress,
		UserAgent:  metadata.UserAgent,
		Referrer:   metadata.Referrer,
		Country:    s.extractCountryFromIP(metadata.IPAddress),
		DeviceType: s.extractDeviceType(metadata.UserAgent),
	}

	// Record the click
	if err := s.analyticsRepo.RecordClick(ctx, click); err != nil {
		return errors.NewDatabaseError(err.Error())
	}

	return nil
}

// GetLinkAnalytics retrieves comprehensive analytics for a specific link
func (s *AnalyticsService) GetLinkAnalytics(ctx context.Context, linkID string) (*models.Analytics, error) {
	if linkID == "" {
		return nil, errors.NewValidationError("Link ID cannot be empty")
	}

	// Verify the link exists
	_, err := s.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.NewLinkNotFoundError(linkID)
		}
		return nil, errors.NewDatabaseError(err.Error())
	}

	// Get analytics summary from repository
	analytics, err := s.analyticsRepo.GetAnalyticsSummary(ctx, linkID)
	if err != nil {
		return nil, errors.NewDatabaseError(err.Error())
	}

	return analytics, nil
}

// GenerateInsights generates AI-powered insights for link analytics
func (s *AnalyticsService) GenerateInsights(ctx context.Context, linkID string) (string, error) {
	if linkID == "" {
		return "", errors.NewValidationError("Link ID cannot be empty")
	}

	// Get analytics data
	analytics, err := s.GetLinkAnalytics(ctx, linkID)
	if err != nil {
		return "", err
	}

	// Generate insights using AI service
	insights, err := s.aiService.AnalyzeClickPatterns(ctx, analytics)
	if err != nil {
		// Return a fallback insight if AI service fails
		return s.generateFallbackInsights(analytics), nil
	}

	return insights, nil
}

// GetClicksByLinkID retrieves all clicks for a specific link
func (s *AnalyticsService) GetClicksByLinkID(ctx context.Context, linkID string) ([]*models.Click, error) {
	if linkID == "" {
		return nil, errors.NewValidationError("Link ID cannot be empty")
	}

	// Verify the link exists
	_, err := s.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errors.NewLinkNotFoundError(linkID)
		}
		return nil, errors.NewDatabaseError(err.Error())
	}

	clicks, err := s.analyticsRepo.GetClicksByLinkID(ctx, linkID)
	if err != nil {
		return nil, errors.NewDatabaseError(err.Error())
	}

	return clicks, nil
}

// GetAnalyticsSummary provides a quick summary of analytics across all links
func (s *AnalyticsService) GetAnalyticsSummary(ctx context.Context) (map[string]interface{}, error) {
	// This could be extended to provide global analytics
	// For now, return a basic structure
	summary := map[string]interface{}{
		"message": "Global analytics summary - to be implemented",
	}

	return summary, nil
}

// Helper methods for metadata extraction

// extractCountryFromIP extracts country code from IP address
// This is a simplified implementation - in production, you'd use a GeoIP service
func (s *AnalyticsService) extractCountryFromIP(ipAddress string) string {
	if ipAddress == "" {
		return "Unknown"
	}

	// Parse IP to validate format
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return "Unknown"
	}

	// Check for local/private IPs
	if ip.IsLoopback() || ip.IsPrivate() {
		return "Local"
	}

	// Simplified country detection - in production, use MaxMind GeoIP or similar
	// For now, return a placeholder
	return "Unknown"
}

// extractDeviceType extracts device type from user agent string
func (s *AnalyticsService) extractDeviceType(userAgent string) string {
	if userAgent == "" {
		return "Unknown"
	}

	userAgent = strings.ToLower(userAgent)

	// Mobile devices
	mobileKeywords := []string{"mobile", "android", "iphone", "ipad", "ipod", "blackberry", "windows phone"}
	for _, keyword := range mobileKeywords {
		if strings.Contains(userAgent, keyword) {
			if strings.Contains(userAgent, "ipad") {
				return "Tablet"
			}
			return "Mobile"
		}
	}

	// Tablet devices
	tabletKeywords := []string{"tablet", "ipad"}
	for _, keyword := range tabletKeywords {
		if strings.Contains(userAgent, keyword) {
			return "Tablet"
		}
	}

	// Bot detection
	botKeywords := []string{"bot", "crawler", "spider", "scraper"}
	for _, keyword := range botKeywords {
		if strings.Contains(userAgent, keyword) {
			return "Bot"
		}
	}

	// Default to desktop
	return "Desktop"
}

// generateFallbackInsights creates basic insights when AI service is unavailable
func (s *AnalyticsService) generateFallbackInsights(analytics *models.Analytics) string {
	if analytics.TotalClicks == 0 {
		return "This link hasn't received any clicks yet."
	}

	insights := fmt.Sprintf("This link has received %d total clicks. ", analytics.TotalClicks)

	// Add device insights
	if len(analytics.DeviceBreakdown) > 0 {
		topDevice := analytics.DeviceBreakdown[0]
		devicePercentage := float64(topDevice.Count) / float64(analytics.TotalClicks) * 100
		insights += fmt.Sprintf("Most clicks (%.1f%%) came from %s devices. ", devicePercentage, strings.ToLower(topDevice.DeviceType))
	}

	// Add country insights
	if len(analytics.CountryBreakdown) > 0 {
		topCountry := analytics.CountryBreakdown[0]
		if topCountry.Country != "Unknown" {
			countryPercentage := float64(topCountry.Count) / float64(analytics.TotalClicks) * 100
			insights += fmt.Sprintf("The top country is %s with %.1f%% of clicks. ", topCountry.Country, countryPercentage)
		}
	}

	// Add daily activity insights
	if len(analytics.DailyClicks) > 0 {
		recentClicks := 0
		for i := 0; i < len(analytics.DailyClicks) && i < 7; i++ {
			recentClicks += analytics.DailyClicks[i].Count
		}
		if recentClicks > 0 {
			insights += fmt.Sprintf("There have been %d clicks in the past week.", recentClicks)
		}
	}

	return insights
}

// GenerateTestClicks creates sample click data for testing purposes
func (s *AnalyticsService) GenerateTestClicks(ctx context.Context, linkID string) error {
	if linkID == "" {
		return errors.NewValidationError("Link ID cannot be empty")
	}

	// Verify the link exists
	_, err := s.linkRepo.GetByID(ctx, linkID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errors.NewLinkNotFoundError(linkID)
		}
		return errors.NewDatabaseError(err.Error())
	}

	// Generate sample clicks
	sampleClicks := s.generateSampleClicks(linkID)
	
	// Insert sample clicks
	for _, click := range sampleClicks {
		if err := s.analyticsRepo.RecordClick(ctx, click); err != nil {
			return errors.NewDatabaseError(err.Error())
		}
	}

	return nil
}

// generateSampleClicks creates realistic sample click data
func (s *AnalyticsService) generateSampleClicks(linkID string) []*models.Click {
	var clicks []*models.Click
	
	countries := []string{"US", "UK", "CA", "DE", "FR", "JP", "AU", "BR"}
	devices := []string{"Desktop", "Mobile", "Tablet"}
	referrers := []string{"", "google.com", "twitter.com", "facebook.com", "linkedin.com"}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		"Mozilla/5.0 (Android 11; Mobile; rv:68.0) Gecko/68.0 Firefox/88.0",
	}
	
	// Generate clicks for the last 7 days
	now := time.Now()
	for i := 0; i < 25; i++ {
		// Random time in the last 7 days
		randomHours := rand.Intn(7 * 24)
		clickTime := now.Add(-time.Duration(randomHours) * time.Hour)
		
		click := &models.Click{
			ID:         uuid.New().String(),
			LinkID:     linkID,
			ClickedAt:  clickTime,
			IPAddress:  fmt.Sprintf("192.168.1.%d", rand.Intn(255)),
			UserAgent:  userAgents[rand.Intn(len(userAgents))],
			Referrer:   referrers[rand.Intn(len(referrers))],
			Country:    countries[rand.Intn(len(countries))],
			DeviceType: devices[rand.Intn(len(devices))],
		}
		
		clicks = append(clicks, click)
	}
	
	return clicks
}