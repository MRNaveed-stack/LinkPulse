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

func TestLinkRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewLinkRepository(mock)
	link := &models.Link{
		ID:             uuid.New(),
		UserID:         uuid.New(),
		Title:          "Git",
		Slug:           "github",
		DestinationURL: "https://github.com",
		IsActive:       true,
		ClickCount:     0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO links`)).
		WithArgs(
			link.ID,
			link.UserID,
			link.Title,
			link.Slug,
			link.DestinationURL,
			link.IsActive,
			link.ClickCount,
			link.CreatedAt,
			link.UpdatedAt,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(link)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestLinkRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewLinkRepository(mock)
	linkID := uuid.New()
	userID := uuid.New()

	rows := pgxmock.NewRows([]string{"id", "user_id", "title", "slug", "destination_url", "is_active", "click_count", "created_at", "updated_at"}).
		AddRow(linkID, userID, "Title", "slug", "https://destination.com", true, int64(0), time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, title, slug, destination_url, is_active, click_count, created_at, updated_at FROM links WHERE id = $1`)).
		WithArgs(linkID).
		WillReturnRows(rows)

	link, err := repo.GetByID(linkID)
	assert.NoError(t, err)
	assert.NotNil(t, link)
	assert.Equal(t, linkID, link.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
