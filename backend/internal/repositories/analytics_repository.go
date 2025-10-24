package repositories

import (
	"context"
	"fmt"

	"url-shortener-backend/internal/database"
	"url-shortener-backend/internal/models"
)

type AnalyticsRepository struct {
	db *database.DB
}

// Ensure AnalyticsRepository implements the interface
var _ AnalyticsRepositoryInterface = (*AnalyticsRepository)(nil)

func NewAnalyticsRepository(db *database.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) RecordClick(ctx context.Context, click *models.Click) error {
	query := `
		INSERT INTO clicks (id, link_id, clicked_at, ip_address, user_agent, referrer, country, device_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Pool.Exec(ctx, query,
		click.ID,
		click.LinkID,
		click.ClickedAt,
		click.IPAddress,
		click.UserAgent,
		click.Referrer,
		click.Country,
		click.DeviceType,
	)

	if err != nil {
		return fmt.Errorf("failed to record click: %w", err)
	}

	return nil
}

func (r *AnalyticsRepository) GetClicksByLinkID(ctx context.Context, linkID string) ([]*models.Click, error) {
	query := `
		SELECT id, link_id, clicked_at, ip_address, user_agent, referrer, country, device_type
		FROM clicks
		WHERE link_id = $1
		ORDER BY clicked_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get clicks by link ID: %w", err)
	}
	defer rows.Close()

	var clicks []*models.Click
	for rows.Next() {
		var click models.Click
		err := rows.Scan(
			&click.ID,
			&click.LinkID,
			&click.ClickedAt,
			&click.IPAddress,
			&click.UserAgent,
			&click.Referrer,
			&click.Country,
			&click.DeviceType,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan click: %w", err)
		}
		clicks = append(clicks, &click)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over clicks: %w", err)
	}

	return clicks, nil
}

func (r *AnalyticsRepository) GetAnalyticsSummary(ctx context.Context, linkID string) (*models.AnalyticsSummary, error) {
	summary := &models.AnalyticsSummary{
		LinkID: linkID,
	}

	// Get total clicks
	totalQuery := `SELECT COUNT(*) FROM clicks WHERE link_id = $1`
	err := r.db.Pool.QueryRow(ctx, totalQuery, linkID).Scan(&summary.TotalClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get total clicks: %w", err)
	}

	// Get unique visitors (unique IP addresses)
	uniqueQuery := `SELECT COUNT(DISTINCT ip_address) FROM clicks WHERE link_id = $1`
	err = r.db.Pool.QueryRow(ctx, uniqueQuery, linkID).Scan(&summary.UniqueVisitors)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique visitors: %w", err)
	}

	return summary, nil

}
