package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/repository"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordResetTokenRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewPasswordResetTokenRepository(mock)
	token := &models.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "mytoken",
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
		CreatedAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO password_reset_tokens`)).
		WithArgs(
			token.ID,
			token.UserID,
			token.Token,
			token.ExpiresAt,
			token.Used,
			token.CreatedAt,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), token)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPasswordResetTokenRepository_GetByToken(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewPasswordResetTokenRepository(mock)
	tokenStr := "some-reset-token-value"
	expectedToken := &models.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     tokenStr,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
		CreatedAt: time.Now(),
	}

	rows := pgxmock.NewRows([]string{"id", "user_id", "token", "expires_at", "used", "created_at"}).
		AddRow(expectedToken.ID, expectedToken.UserID, expectedToken.Token, expectedToken.ExpiresAt, expectedToken.Used, expectedToken.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, token, expires_at, used, created_at FROM password_reset_tokens WHERE token = $1`)).
		WithArgs(tokenStr).
		WillReturnRows(rows)

	token, err := repo.GetByToken(context.Background(), tokenStr)
	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.Equal(t, tokenStr, token.Token)
	assert.NoError(t, mock.ExpectationsWereMet())
}
