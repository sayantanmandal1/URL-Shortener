package testutils

import (
	"context"

	"url-shortener-backend/internal/models"
	"url-shortener-backend/internal/repositories"

	"github.com/stretchr/testify/assert"
)

// NotFoundError represents a not found error
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// MockLinkRepository is a mock implementation for testing
type MockLinkRepository struct {
	links map[string]*models.Link
	slugs map[string]*models.Link
}

// Ensure MockLinkRepository implements the interface
var _ repositories.LinkRepositoryInterface = (*MockLinkRepository)(nil)

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
	return nil, &NotFoundError{Message: "link not found"} // Simulate not found error
}

func (m *MockLinkRepository) GetBySlug(ctx context.Context, slug string) (*models.Link, error) {
	if link, exists := m.slugs[slug]; exists {
		return link, nil
	}
	return nil, &NotFoundError{Message: "link not found"} // Simulate not found error
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

// MockAnalyticsRepository is a mock implementation for testing
type MockAnalyticsRepository struct {
	clicks map[string][]*models.Click
}

// Ensure MockAnalyticsRepository implements the interface
var _ repositories.AnalyticsRepositoryInterface = (*MockAnalyticsRepository)(nil)

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

// MockAIService is a mock implementation for testing
type MockAIService struct{}

func NewMockAIService() *MockAIService {
	return &MockAIService{}
}

func (m *MockAIService) GenerateTitleAndDescription(ctx context.Context, url string) (*models.AIMetadata, error) {
	return &models.AIMetadata{
		Title:       "Mock Title",
		Description: "Mock Description",
	}, nil
}

func (m *MockAIService) AnalyzeClickPatterns(ctx context.Context, analytics *models.AnalyticsSummary) (string, error) {
	return "Mock AI insights", nil
}
