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

func TestUserRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserRepository(mock)
	user := &models.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@linkpulse.com",
		PasswordHash: "hash123",
		Plan:         "free",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users`)).
		WithArgs(
			user.ID,
			user.Username,
			user.Email,
			user.PasswordHash,
			user.Plan,
			user.GoogleID,
			pgxmock.AnyArg(), // CreatedAt
			pgxmock.AnyArg(), // UpdatedAt
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewUserRepository(mock)
	email := "test@linkpulse.com"
	expectedUser := &models.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        email,
		PasswordHash: "hash123",
		Plan:         "free",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	rows := pgxmock.NewRows([]string{"id", "username", "email", "password_hash", "plan", "google_id", "created_at", "updated_at"}).
		AddRow(expectedUser.ID, expectedUser.Username, expectedUser.Email, expectedUser.PasswordHash, expectedUser.Plan, expectedUser.GoogleID, expectedUser.CreatedAt, expectedUser.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password_hash, plan, google_id, created_at, updated_at FROM users WHERE email = $1`)).
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.GetByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
