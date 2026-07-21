package models

import (
	"time"

	"github.com/google/uuid"
)

type ClickEvent struct {
	ID     uuid.UUID
	LinkID uuid.UUID

	IPAddress string
	UserAgent string
	Referrer  string

	ClickedAt time.Time
}
