package testutils

import (
	"context"

	"url-shortener-backend/internal/models"
)

// MockAIService is a mock implementation of AIService for testing
type MockAIService struct {
	GenerateTitleAndDescriptionFunc func(ctx context.Context, url string) (*models.AIMetadata, error)
	AnalyzeClickPatternsFunc        func(ctx context.Context, analytics *models.Analytics) (string, error)
}

func (m *MockAIService) GenerateTitleAndDescription(ctx context.Context, url string) (*models.AIMetadata, error) {
	if m.GenerateTitleAndDescriptionFunc != nil {
		return m.GenerateTitleAndDescriptionFunc(ctx, url)
	}
	return &models.AIMetadata{
		Title:       "Mock Title",
		Description: "Mock Description",
	}, nil
}

func (m *MockAIService) AnalyzeClickPatterns(ctx context.Context, analytics *models.Analytics) (string, error) {
	if m.AnalyzeClickPatternsFunc != nil {
		return m.AnalyzeClickPatternsFunc(ctx, analytics)
	}
	return "Mock insights: This link is performing well.", nil
}

// NewMockAIService creates a new mock AI service
func NewMockAIService() *MockAIService {
	return &MockAIService{}
}
