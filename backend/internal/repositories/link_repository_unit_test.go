package repositories

import (
	"context"
	"testing"
	"time"

	"url-shortener-backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Simple unit tests without database dependencies

func TestLinkRepositoryInterface(t *testing.T) {
	// Test that our concrete types implement the interfaces
	var _ LinkRepositoryInterface = (*LinkRepository)(nil)
	var _ AnalyticsRepositoryInterface = (*AnalyticsRepository)(nil)
}

func TestLinkModel(t *testing.T) {
	// Test link model creation and validation
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

	assert.NotEmpty(t, link.ID)
	assert.Equal(t, "https://example.com", link.OriginalURL)
	assert.Equal(t, "test-slug", link.Slug)
	assert.Equal(t, "Test Title", link.Title)
	assert.Equal(t, "Test Description", link.Description)
	assert.Equal(t, 0, link.ClickCount)
	assert.False(t, link.CreatedAt.IsZero())
	assert.False(t, link.UpdatedAt.IsZero())
}

func TestClickModel(t *testing.T) {
	// Test click model creation and validation
	click := &models.Click{
		ID:         uuid.New().String(),
		LinkID:     uuid.New().String(),
		ClickedAt:  time.Now(),
		IPAddress:  "192.168.1.1",
		UserAgent:  "Mozilla/5.0 Test Browser",
		Referrer:   "https://google.com",
		Country:    "US",
		DeviceType: "Desktop",
	}

	assert.NotEmpty(t, click.ID)
	assert.NotEmpty(t, click.LinkID)
	assert.False(t, click.ClickedAt.IsZero())
	assert.Equal(t, "192.168.1.1", click.IPAddress)
	assert.Equal(t, "Mozilla/5.0 Test Browser", click.UserAgent)
	assert.Equal(t, "https://google.com", click.Referrer)
	assert.Equal(t, "US", click.Country)
	assert.Equal(t, "Desktop", click.DeviceType)
}

func TestAnalyticsSummaryModel(t *testing.T) {
	// Test analytics summary model
	summary := &models.AnalyticsSummary{
		LinkID:         uuid.New().String(),
		TotalClicks:    100,
		UniqueVisitors: 75,
	}

	assert.NotEmpty(t, summary.LinkID)
	assert.Equal(t, 100, summary.TotalClicks)
	assert.Equal(t, 75, summary.UniqueVisitors)
}

// Mock repository for testing service layer
type MockLinkRepository struct {
	links map[string]*models.Link
	slugs map[string]*models.Link
}

func NewMockLinkRepository() *MockLinkRepository {
	return &MockLinkRepository{
		links: make(map[string]*models.Link),
		slugs: make(map[string]*models.Link),
	}
}

func (m *MockLinkRepository) Create(ctx context.Context, link *models.Link) error {
	if _, exists := m.slugs[link.Slug]; exists {
		return assert.AnError // Simulate duplicate slug error
	}
	m.links[link.ID] = link
	m.slugs[link.Slug] = link
	return nil
}

func (m *MockLinkRepository) GetByID(ctx context.Context, id string) (*models.Link, error) {
	if link, exists := m.links[id]; exists {
		return link, nil
	}
	return nil, assert.AnError // Simulate not found error
}

func (m *MockLinkRepository) GetBySlug(ctx context.Context, slug string) (*models.Link, error) {
	if link, exists := m.slugs[slug]; exists {
		return link, nil
	}
	return nil, assert.AnError // Simulate not found error
}

func (m *MockLinkRepository) GetAll(ctx context.Context) ([]*models.Link, error) {
	var links []*models.Link
	for _, link := range m.links {
		links = append(links, link)
	}
	return links, nil
}

func (m *MockLinkRepository) IncrementClicks(ctx context.Context, id string) error {
	if link, exists := m.links[id]; exists {
		link.ClickCount++
		return nil
	}
	return assert.AnError // Simulate not found error
}

// Mock analytics repository
type MockAnalyticsRepository struct {
	clicks map[string][]*models.Click
}

func NewMockAnalyticsRepository() *MockAnalyticsRepository {
	return &MockAnalyticsRepository{
		clicks: make(map[string][]*models.Click),
	}
}

func (m *MockAnalyticsRepository) RecordClick(ctx context.Context, click *models.Click) error {
	m.clicks[click.LinkID] = append(m.clicks[click.LinkID], click)
	return nil
}

func (m *MockAnalyticsRepository) GetClicksByLinkID(ctx context.Context, linkID string) ([]*models.Click, error) {
	return m.clicks[linkID], nil
}

func (m *MockAnalyticsRepository) GetAnalyticsSummary(ctx context.Context, linkID string) (*models.AnalyticsSummary, error) {
	clicks := m.clicks[linkID]
	uniqueIPs := make(map[string]bool)

	for _, click := range clicks {
		uniqueIPs[click.IPAddress] = true
	}

	return &models.AnalyticsSummary{
		LinkID:         linkID,
		TotalClicks:    len(clicks),
		UniqueVisitors: len(uniqueIPs),
	}, nil
}

func TestMockRepositories(t *testing.T) {
	// Test mock link repository
	linkRepo := NewMockLinkRepository()
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

	// Test create
	err := linkRepo.Create(ctx, link)
	assert.NoError(t, err)

	// Test get by ID
	retrieved, err := linkRepo.GetByID(ctx, link.ID)
	assert.NoError(t, err)
	assert.Equal(t, link.ID, retrieved.ID)

	// Test get by slug
	retrieved, err = linkRepo.GetBySlug(ctx, link.Slug)
	assert.NoError(t, err)
	assert.Equal(t, link.Slug, retrieved.Slug)

	// Test increment clicks
	err = linkRepo.IncrementClicks(ctx, link.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, link.ClickCount)

	// Test analytics repository
	analyticsRepo := NewMockAnalyticsRepository()

	click := &models.Click{
		ID:         uuid.New().String(),
		LinkID:     link.ID,
		ClickedAt:  time.Now(),
		IPAddress:  "192.168.1.1",
		UserAgent:  "Test Browser",
		Referrer:   "",
		Country:    "US",
		DeviceType: "Desktop",
	}

	err = analyticsRepo.RecordClick(ctx, click)
	assert.NoError(t, err)

	clicks, err := analyticsRepo.GetClicksByLinkID(ctx, link.ID)
	assert.NoError(t, err)
	assert.Len(t, clicks, 1)

	summary, err := analyticsRepo.GetAnalyticsSummary(ctx, link.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, summary.TotalClicks)
	assert.Equal(t, 1, summary.UniqueVisitors)
}
