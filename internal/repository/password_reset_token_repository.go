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

type PasswordResetTokenRepository interface {
	Create(ctx context.Context, token *models.PasswordResetToken) error
	GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error)
	MarkUsed(ctx context.Context, tokenID uuid.UUID, usedAt time.Time) error
}

type passwordResetTokenRepository struct {
	db DB
}

func NewPasswordResetTokenRepository(db DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

var (
	ErrResetTokenNotFound = errors.New("password reset token not found")
)

func (r *passwordResetTokenRepository) Create(ctx context.Context, t *models.PasswordResetToken) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.ExpiresAt.IsZero() {
		return fmt.Errorf("expires_at is required")
	}

	query := `
		INSERT INTO password_reset_tokens (id, user_id, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		t.ID,
		t.UserID,
		t.Token,
		t.ExpiresAt,
		t.Used,
		t.CreatedAt,
	)
	if err != nil {
		logger.Log.Error(
			"failed to create password reset token",
			"error", err,
		)
		return fmt.Errorf("failed to create password reset token: %w", err)
	}
	return nil
}

func (r *passwordResetTokenRepository) GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token = $1
		LIMIT 1
	`

	var t models.PasswordResetToken
	err := r.db.QueryRow(ctx, query, token).Scan(
		&t.ID,
		&t.UserID,
		&t.Token,
		&t.ExpiresAt,
		&t.Used,
		&t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrResetTokenNotFound
		}
		logger.Log.Error(
			"failed to get reset token",
			"error", err,
		)
		return nil, fmt.Errorf("failed to get reset token: %w", err)
	}

	return &t, nil
}

func (r *passwordResetTokenRepository) MarkUsed(ctx context.Context, tokenID uuid.UUID, usedAt time.Time) error {
	query := `
		UPDATE password_reset_tokens
		SET used = TRUE
		WHERE id = $1
	`

	res, err := r.db.Exec(ctx, query, tokenID)
	_ = usedAt // reserved if you later add used_at column
	if err != nil {
		logger.Log.Error(
			"failed to mark reset token used",
			"error", err,
		)
		return fmt.Errorf("failed to mark reset token used: %w", err)
	}
	_ = res
	return nil
}
