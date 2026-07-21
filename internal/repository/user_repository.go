package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/MRNaveed-stack/LinkPulse/internal/logger"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*models.User, error)
	LinkGoogleAccount(ctx context.Context, id uuid.UUID, googleID string) error
	UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error
	UpdatePlan(ctx context.Context, id uuid.UUID, plan string) error
	UpdateUsername(ctx context.Context, id uuid.UUID, username string) error
	UpdateEmail(ctx context.Context, id uuid.UUID, email string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db DB
}

func NewUserRepository(db DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, plan, google_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	if user.Plan == "" {
		user.Plan = "free"
	}

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Plan,
		user.GoogleID,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		logger.Log.Error(
			"failed to create user",
			"error", err,
		)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, plan, google_id, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Plan,
		&user.GoogleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Log.Error(
			"failed to get user by email",
			"error", err,
		)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, plan, google_id, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Plan,
		&user.GoogleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Log.Error(
			"failed to get user by id",
			"error", err,
		)
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, plan, google_id, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Plan,
		&user.GoogleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Log.Error(
			"failed to get user by username",
			"error", err,
		)
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, hash, time.Now(), id)
	if err != nil {
		logger.Log.Error(
			"failed to update password",
			"error", err,
		)
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (r *userRepository) UpdatePlan(ctx context.Context, id uuid.UUID, plan string) error {
	query := `
		UPDATE users
		SET plan = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, plan, time.Now(), id)
	if err != nil {
		logger.Log.Error(
			"failed to update plan",
			"error", err,
		)
		return fmt.Errorf("failed to update plan: %w", err)
	}

	return nil
}

func (r *userRepository) GetByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, plan, google_id, created_at, updated_at
		FROM users
		WHERE google_id = $1
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, googleID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Plan,
		&user.GoogleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logger.Log.Error(
			"failed to get user by google id",
			"error", err,
		)
		return nil, fmt.Errorf("failed to get user by google id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) LinkGoogleAccount(ctx context.Context, id uuid.UUID, googleID string) error {
	query := `
		UPDATE users
		SET google_id = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, googleID, time.Now(), id)
	if err != nil {
		logger.Log.Error(
			"failed to link google account",
			"error", err,
		)
		return fmt.Errorf("failed to link google account: %w", err)
	}

	return nil
}

func (r *userRepository) UpdateUsername(ctx context.Context, id uuid.UUID, username string) error {
	query := `
		UPDATE users
		SET username = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, username, time.Now(), id)
	if err != nil {
		logger.Log.Error("failed to update username", "error", err)
		return fmt.Errorf("failed to update username: %w", err)
	}
	return nil
}

func (r *userRepository) UpdateEmail(ctx context.Context, id uuid.UUID, email string) error {
	query := `
		UPDATE users
		SET email = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(ctx, query, email, time.Now(), id)
	if err != nil {
		logger.Log.Error("failed to update email", "error", err)
		return fmt.Errorf("failed to update email: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.Log.Error("failed to delete user", "error", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
