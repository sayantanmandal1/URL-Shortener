package repositories

import (
	"context"
	"fmt"

	"url-shortener-backend/internal/database"
	"url-shortener-backend/internal/models"

	"github.com/jackc/pgx/v5"
)

type LinkRepository struct {
	db *database.DB
}

func NewLinkRepository(db *database.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) Create(ctx context.Context, link *models.Link) error {
	query := `
		INSERT INTO links (id, original_url, slug, title, description, click_count, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Pool.Exec(ctx, query,
		link.ID,
		link.OriginalURL,
		link.Slug,
		link.Title,
		link.Description,
		link.ClickCount,
		link.CreatedAt,
		link.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create link: %w", err)
	}

	return nil
}

func (r *LinkRepository) GetByID(ctx context.Context, id string) (*models.Link, error) {
	query := `
		SELECT id, original_url, slug, title, description, click_count, created_at, updated_at
		FROM links
		WHERE id = $1`

	var link models.Link
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&link.ID,
		&link.OriginalURL,
		&link.Slug,
		&link.Title,
		&link.Description,
		&link.ClickCount,
		&link.CreatedAt,
		&link.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("link not found")
		}
		return nil, fmt.Errorf("failed to get link by ID: %w", err)
	}

	return &link, nil
}

func (r *LinkRepository) GetBySlug(ctx context.Context, slug string) (*models.Link, error) {
	query := `
		SELECT id, original_url, slug, title, description, click_count, created_at, updated_at
		FROM links
		WHERE slug = $1`

	var link models.Link
	err := r.db.Pool.QueryRow(ctx, query, slug).Scan(
		&link.ID,
		&link.OriginalURL,
		&link.Slug,
		&link.Title,
		&link.Description,
		&link.ClickCount,
		&link.CreatedAt,
		&link.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("link not found")
		}
		return nil, fmt.Errorf("failed to get link by slug: %w", err)
	}

	return &link, nil
}

func (r *LinkRepository) GetAll(ctx context.Context) ([]*models.Link, error) {
	query := `
		SELECT id, original_url, slug, title, description, click_count, created_at, updated_at
		FROM links
		ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all links: %w", err)
	}
	defer rows.Close()

	var links []*models.Link
	for rows.Next() {
		var link models.Link
		err := rows.Scan(
			&link.ID,
			&link.OriginalURL,
			&link.Slug,
			&link.Title,
			&link.Description,
			&link.ClickCount,
			&link.CreatedAt,
			&link.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan link: %w", err)
		}
		links = append(links, &link)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over links: %w", err)
	}

	return links, nil
}

func (r *LinkRepository) IncrementClicks(ctx context.Context, id string) error {
	query := `
		UPDATE links 
		SET click_count = click_count + 1, updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment clicks: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("link not found")
	}

	return nil
}
