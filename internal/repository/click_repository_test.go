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

func TestClickRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewClickRepository(mock)
	event := &models.ClickEvent{
		ID:        uuid.New(),
		LinkID:    uuid.New(),
		UserAgent: "Mozilla/5.0",
		IPAddress: "192.168.1.1",
		Referrer:  "https://google.com",
		ClickedAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO click_events`)).
		WithArgs(
			event.ID,
			event.LinkID,
			event.UserAgent,
			event.IPAddress,
			event.Referrer,
			event.ClickedAt,
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(event)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
