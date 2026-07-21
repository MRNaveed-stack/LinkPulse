package repository_test

import (
	"regexp"
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/repository"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkRepository_GetAnalyticsOverview(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := repository.NewLinkRepository(mock)
	userID := uuid.New()

	// 1. Mock overview counters query
	overviewRows := pgxmock.NewRows([]string{"total_links", "total_clicks", "clicks_today", "active_links"}).
		AddRow(int64(5), int64(150), int64(10), int64(4))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT
		(SELECT COUNT(*) FROM links WHERE user_id = $1) as total_links,
		(SELECT COALESCE(SUM(click_count), 0) FROM links WHERE user_id = $1) as total_clicks,
		(SELECT COUNT(*) FROM click_events ce JOIN links l ON ce.link_id = l.id WHERE l.user_id = $1 AND ce.clicked_at >= CURRENT_DATE) as clicks_today,
		(SELECT COUNT(*) FROM links WHERE user_id = $1 AND is_active = true) as active_links`)).
		WithArgs(userID).
		WillReturnRows(overviewRows)

	// 2. Mock top link query
	topLinkRows := pgxmock.NewRows([]string{"title", "slug", "click_count"}).
		AddRow("Google link", "google", int64(7))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT l.title, l.slug, COUNT(ce.id) AS click_count
	FROM links l
	JOIN click_events ce ON l.id = ce.link_id
	WHERE l.user_id = $1
	GROUP BY l.id, l.title, l.slug
	ORDER BY click_count DESC, l.created_at DESC
	LIMIT 1`)).
		WithArgs(userID).
		WillReturnRows(topLinkRows)

	overview, err := repo.GetAnalyticsOverview(userID)
	assert.NoError(t, err)
	assert.NotNil(t, overview)

	assert.Equal(t, int64(5), overview.TotalLinks)
	assert.Equal(t, int64(150), overview.TotalClicks)
	assert.Equal(t, int64(10), overview.ClicksToday)
	assert.Equal(t, int64(4), overview.ActiveLinks)

	assert.NotNil(t, overview.TopLink)
	assert.Equal(t, "Google link", overview.TopLink.Title)
	assert.Equal(t, "google", overview.TopLink.Slug)
	assert.Equal(t, int64(7), overview.TopLink.Clicks)

	assert.NoError(t, mock.ExpectationsWereMet())
}
