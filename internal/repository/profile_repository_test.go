package repository_test

import (
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

func TestProfileRepository_GetByUserID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewProfileRepository(mock)
	userID := uuid.New()

	rows := pgxmock.NewRows([]string{"user_id", "display_name", "bio", "avatar_url", "created_at", "updated_at"}).
		AddRow(userID, "Disp", "Bio", "url", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT user_id, display_name, bio, avatar_url, created_at, updated_at FROM profiles WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(rows)

	profile, err := repo.GetByUserID(userID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, userID, profile.UserID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProfileRepository_UpdateProfile(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewProfileRepository(mock)
	profile := &models.Profile{
		UserID:      uuid.New(),
		DisplayName: "New Name",
		Bio:         "New Bio",
		AvatarURL:   "http://newavatar.url",
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE profiles SET display_name = $1, bio = $2, avatar_url = $3`)).
		WithArgs(profile.DisplayName, profile.Bio, profile.AvatarURL, profile.UserID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateProfile(profile)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProfileRepository_CreateProfile(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewProfileRepository(mock)
	profile := &models.Profile{
		UserID:      uuid.New(),
		DisplayName: "Name",
		Bio:         "Bio",
		AvatarURL:   "http://avatar.url",
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO profiles`)).
		WithArgs(profile.UserID, profile.DisplayName, profile.Bio, profile.AvatarURL).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.CreateProfile(profile)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
