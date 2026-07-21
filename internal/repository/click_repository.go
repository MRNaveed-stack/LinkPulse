package repository

import (
	"context"

	"github.com/MRNaveed-stack/LinkPulse/internal/logger"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
)

type ClickRepository interface {
	Create(event *models.ClickEvent) error
}

type clickRepository struct {
	db DB
}

func NewClickRepository(db DB) ClickRepository {
	return &clickRepository{db: db}
}

func (r *clickRepository) Create(event *models.ClickEvent) error {
	query := `
	INSERT INTO click_events (id, link_id, user_agent, ip_address, 
	referrer , clicked_at
	) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(
		context.Background(),
		query,
		event.ID,
		event.LinkID,
		event.UserAgent,
		event.IPAddress,
		event.Referrer,
		event.ClickedAt,
	)
	if err != nil {
		logger.Log.Error(
			"failed to create click event",
			"error", err,
		)
	}
	return err
}
