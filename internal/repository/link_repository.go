package repository

import (
	"context"
	"errors"

	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type LinkRepository interface {
	Create(link *models.Link) error
	GetByID(id uuid.UUID) (*models.Link, error)
	GetBySlug(slug string) (*models.Link, error)
	GetByUserIDAndSlug(userID uuid.UUID, slug string) (*models.Link, error)
	GetByUserID(userID uuid.UUID) ([]*models.Link, error)
	GetActiveLinksByUserID(userID uuid.UUID) ([]*models.Link, error)
	IncrementClickCount(id uuid.UUID) error
	Update(link *models.Link) error
	Delete(id uuid.UUID) error
	UpdateStatus(id uuid.UUID, isActive bool) error
	GetAnalyticsOverview(userID uuid.UUID) (*models.AnalyticsOverview, error)
	GetLinkAnalytics(userID uuid.UUID) ([]models.LinkAnalyticsDTO, error)
	GetDailyAnalytics(userID uuid.UUID, days int) ([]models.DailyAnalyticsDTO, error)
	GetRecentActivity(userID uuid.UUID, limit int) ([]models.RecentActivityDTO, error)
	GetReferrerAnalytics(userID uuid.UUID) ([]models.ReferrerAnalyticsDTO, error)
}

type linkRepository struct {
	db DB
}

func NewLinkRepository(db DB) LinkRepository {
	return &linkRepository{db: db}
}

func (r *linkRepository) Create(link *models.Link) error {
	query := `
	INSERT INTO links (id , user_id , title, slug, destination_url, is_active,click_count,created_at,updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		link.ID,
		link.UserID,
		link.Title,
		link.Slug,
		link.DestinationURL,
		link.IsActive,
		link.ClickCount,
		link.CreatedAt,
		link.UpdatedAt,
	)
	return err
}

func (r *linkRepository) GetByID(id uuid.UUID) (*models.Link, error) {
	query := `
	SELECT id, user_id, title, slug, destination_url, is_active, click_count, created_at, updated_at
	FROM links WHERE id = $1
	`
	var link models.Link
	err := r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&link.ID,
		&link.UserID,
		&link.Title,
		&link.Slug,
		&link.DestinationURL,
		&link.IsActive,
		&link.ClickCount,
		&link.CreatedAt,
		&link.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &link, nil
}

func (r *linkRepository) GetByUserID(userID uuid.UUID) ([]*models.Link, error) {
	query := `
	SELECT id, user_id, title , slug, destination_url , is_active, click_count, created_at , updated_at FROM links WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(
		context.Background(),
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var links []*models.Link
	for rows.Next() {
		var link models.Link
		err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.Title,
			&link.Slug,
			&link.DestinationURL,
			&link.IsActive,
			&link.ClickCount,
			&link.CreatedAt,
			&link.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	return links, nil
}

func (r *linkRepository) GetBySlug(slug string) (*models.Link, error) {
	query := `
	SELECT id, user_id, title , slug, destination_url, is_active,
	click_count, created_at, updated_at
	FROM links WHERE slug = $1
	`

	var link models.Link
	err := r.db.QueryRow(
		context.Background(),
		query,
		slug,
	).Scan(
		&link.ID,
		&link.UserID,
		&link.Title,
		&link.Slug,
		&link.DestinationURL,
		&link.IsActive,
		&link.ClickCount,
		&link.CreatedAt,
		&link.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &link, nil
}

func (r *linkRepository) GetByUserIDAndSlug(userID uuid.UUID, slug string) (*models.Link, error) {
	query := `
	SELECT id, user_id, title , slug, destination_url, is_active,
	click_count, created_at, updated_at
	FROM links WHERE user_id = $1 AND slug = $2
	`

	var link models.Link
	err := r.db.QueryRow(
		context.Background(),
		query,
		userID,
		slug,
	).Scan(
		&link.ID,
		&link.UserID,
		&link.Title,
		&link.Slug,
		&link.DestinationURL,
		&link.IsActive,
		&link.ClickCount,
		&link.CreatedAt,
		&link.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &link, nil
}

func (r *linkRepository) IncrementClickCount(
	linkID uuid.UUID,
) error {
	query := `
	UPDATE links
	SET click_count = click_count + 1
	WHERE id = $1`
	_, err := r.db.Exec(
		context.Background(),
		query,
		linkID,
	)
	return err
}

func (r *linkRepository) GetActiveLinksByUserID(
	userID uuid.UUID,
) ([]*models.Link, error) {
	query := `
	SELECT id, user_id, 
	title , slug, destination_url, 
	is_active, click_count, created_at, updated_at
	FROM links WHERE user_id = $1 AND is_active = true
	ORDER BY created_at DESC
	`
	rows, err := r.db.Query(
		context.Background(),
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*models.Link
	for rows.Next() {
		var link models.Link
		err := rows.Scan(
			&link.ID,
			&link.UserID,
			&link.Title,
			&link.Slug,
			&link.DestinationURL,
			&link.IsActive,
			&link.ClickCount,
			&link.CreatedAt,
			&link.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, &link)
	}
	return links, nil
}

func (r *linkRepository) Update(link *models.Link) error {
	query := `
	UPDATE links 
	SET
	title = $1,
	slug = $2,
	destination_url = $3,
	updated_at = NOW()
	WHERE id = $4
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		link.Title,
		link.Slug,
		link.DestinationURL,
		link.ID,
	)
	return err
}

func (r *linkRepository) Delete(id uuid.UUID) error {
	query := `
	DELETE FROM links WHERE id = $1
	`
	_, err := r.db.Exec(context.Background(), query, id)
	return err
}

func (r *linkRepository) UpdateStatus(id uuid.UUID, isActive bool) error {
	query := `
	UPDATE links
	SET is_active = $1, updated_at = NOW()
	WHERE id = $2
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		isActive,
		id,
	)
	return err
}
