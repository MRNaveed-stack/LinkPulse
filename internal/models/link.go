package models

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`

	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	DestinationURL string    `json:"destination_url"`
	IsActive       bool      `json:"is_active"`
	ClickCount     int64     `json:"click_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UpdateLinkStatusRequest struct {
	IsActive bool `json:"is_active"`
}
