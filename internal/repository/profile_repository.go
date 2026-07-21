package repository

import (
	"context"
	"errors"

	"github.com/MRNaveed-stack/LinkPulse/internal/logger"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProfileRepository interface {
	GetByUserID(userID uuid.UUID) (*models.Profile, error)
	CreateProfile(profile *models.Profile) error
	UpdateProfile(profile *models.Profile) error
}

type profileRepository struct {
	db DB
}

func NewProfileRepository(db DB) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) GetByUserID(userID uuid.UUID) (*models.Profile, error) {
	query := `
   SELECT user_id, display_name, bio, avatar_url, created_at, updated_at	
   FROM profiles WHERE user_id = $1
	`
	var profile models.Profile
	err := r.db.QueryRow(
		context.Background(),
		query,
		userID,
	).Scan(
		&profile.UserID,
		&profile.DisplayName,
		&profile.Bio,
		&profile.AvatarURL,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Log.Error(
			"failed to get profile by user id",
			"error", err,
		)
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) UpdateProfile(profile *models.Profile) error {
	query := `
	UPDATE profiles SET
	display_name = $1, bio = $2, avatar_url = $3, updated_at = NOW() WHERE user_id = $4
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		profile.DisplayName,
		profile.Bio,
		profile.AvatarURL,
		profile.UserID,
	)
	if err != nil {
		logger.Log.Error(
			"failed to update profile",
			"error", err,
		)
	}
	return err
}

func (r *profileRepository) CreateProfile(profile *models.Profile) error {
	query := `
	INSERT INTO profiles (user_id, display_name, bio, avatar_url, created_at, updated_at)
	VALUES ($1, $2,$3,$4,NOW(), NOW())
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		profile.UserID,
		profile.DisplayName,
		profile.Bio,
		profile.AvatarURL,
	)
	if err != nil {
		logger.Log.Error(
			"failed to create profile",
			"error", err,
		)
	}
	return err
}
