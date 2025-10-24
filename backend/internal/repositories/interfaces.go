package repositories

import (
	"context"
	"url-shortener-backend/internal/models"
)

// LinkRepositoryInterface defines the interface for link repository operations
type LinkRepositoryInterface interface {
	Create(ctx context.Context, link *models.Link) error
	GetByID(ctx context.Context, id string) (*models.Link, error)
	GetBySlug(ctx context.Context, slug string) (*models.Link, error)
	GetAll(ctx context.Context) ([]*models.Link, error)
	IncrementClicks(ctx context.Context, id string) error
}

// AnalyticsRepositoryInterface defines the interface for analytics repository operations
type AnalyticsRepositoryInterface interface {
	RecordClick(ctx context.Context, click *models.Click) error
	GetClicksByLinkID(ctx context.Context, linkID string) ([]*models.Click, error)
	GetAnalyticsSummary(ctx context.Context, linkID string) (*models.AnalyticsSummary, error)
}
