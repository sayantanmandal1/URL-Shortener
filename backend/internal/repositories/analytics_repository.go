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

func (r *AnalyticsRepository) GetAnalyticsSummary(ctx context.Context, linkID string) (*models.Analytics, error) {
	analytics := &models.Analytics{
		DailyClicks:       []models.DailyClick{},
		CountryBreakdown:  []models.CountryStats{},
		DeviceBreakdown:   []models.DeviceStats{},
		ReferrerBreakdown: []models.ReferrerStats{},
	}

	// Get total clicks
	totalQuery := `SELECT COUNT(*) FROM clicks WHERE link_id = $1`
	err := r.db.Pool.QueryRow(ctx, totalQuery, linkID).Scan(&analytics.TotalClicks)
	if err != nil {
		return nil, fmt.Errorf("failed to get total clicks: %w", err)
	}

	// Get daily clicks for the last 30 days
	dailyQuery := `
		SELECT DATE(clicked_at) as date, COUNT(*) as count
		FROM clicks
		WHERE link_id = $1 AND clicked_at >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(clicked_at)
		ORDER BY date DESC`

	rows, err := r.db.Pool.Query(ctx, dailyQuery, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily clicks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var daily models.DailyClick
		err := rows.Scan(&daily.Date, &daily.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan daily click: %w", err)
		}
		analytics.DailyClicks = append(analytics.DailyClicks, daily)
	}

	// Get country breakdown
	countryQuery := `
		SELECT COALESCE(country, 'Unknown') as country, COUNT(*) as count
		FROM clicks
		WHERE link_id = $1
		GROUP BY country
		ORDER BY count DESC
		LIMIT 10`

	rows, err = r.db.Pool.Query(ctx, countryQuery, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get country breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var country models.CountryStats
		err := rows.Scan(&country.Country, &country.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan country stats: %w", err)
		}
		analytics.CountryBreakdown = append(analytics.CountryBreakdown, country)
	}

	// Get device breakdown
	deviceQuery := `
		SELECT COALESCE(device_type, 'Unknown') as device_type, COUNT(*) as count
		FROM clicks
		WHERE link_id = $1
		GROUP BY device_type
		ORDER BY count DESC`

	rows, err = r.db.Pool.Query(ctx, deviceQuery, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var device models.DeviceStats
		err := rows.Scan(&device.DeviceType, &device.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan device stats: %w", err)
		}
		analytics.DeviceBreakdown = append(analytics.DeviceBreakdown, device)
	}

	// Get referrer breakdown
	referrerQuery := `
		SELECT COALESCE(referrer, 'Direct') as referrer, COUNT(*) as count
		FROM clicks
		WHERE link_id = $1
		GROUP BY referrer
		ORDER BY count DESC
		LIMIT 10`

	rows, err = r.db.Pool.Query(ctx, referrerQuery, linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get referrer breakdown: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var referrer models.ReferrerStats
		err := rows.Scan(&referrer.Referrer, &referrer.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan referrer stats: %w", err)
		}
		analytics.ReferrerBreakdown = append(analytics.ReferrerBreakdown, referrer)
	}

	return analytics, nil
}
