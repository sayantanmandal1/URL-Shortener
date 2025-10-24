package models

import (
	"time"
)

type Link struct {
	ID          string    `json:"id" db:"id"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	Slug        string    `json:"slug" db:"slug"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	ClickCount  int       `json:"click_count" db:"click_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url" validate:"required,url"`
	CustomSlug  string `json:"custom_slug,omitempty" validate:"omitempty,alphanum,min=3,max=50"`
}

// AIMetadata represents AI-generated metadata for a link
type AIMetadata struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
