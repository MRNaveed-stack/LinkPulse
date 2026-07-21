package repository

import (
	"context"
	"errors"
	"sort"

	"github.com/MRNaveed-stack/LinkPulse/internal/logger"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *linkRepository) GetAnalyticsOverview(
	userID uuid.UUID,
) (*models.AnalyticsOverview, error) {
	query := `
	SELECT
		(SELECT COUNT(*) FROM links WHERE user_id = $1) as total_links,
		(SELECT COALESCE(SUM(click_count), 0) FROM links WHERE user_id = $1) as total_clicks,
		(SELECT COUNT(*) FROM click_events ce JOIN links l ON ce.link_id = l.id WHERE l.user_id = $1 AND ce.clicked_at >= CURRENT_DATE) as clicks_today,
		(SELECT COUNT(*) FROM links WHERE user_id = $1 AND is_active = true) as active_links
	`
	var overview models.AnalyticsOverview
	err := r.db.QueryRow(context.Background(),
		query,
		userID,
	).Scan(
		&overview.TotalLinks,
		&overview.TotalClicks,
		&overview.ClicksToday,
		&overview.ActiveLinks,
	)
	if err != nil {
		logger.Log.Error(
			"failed to scan analytics overview",
			"error", err,
		)
		return nil, err
	}

	topLinkQuery := `
	SELECT l.title, l.slug, COUNT(ce.id) AS click_count
	FROM links l
	JOIN click_events ce ON l.id = ce.link_id
	WHERE l.user_id = $1
	GROUP BY l.id, l.title, l.slug
	ORDER BY click_count DESC, l.created_at DESC
	LIMIT 1
	`
	var topLinkTitle string
	var topLinkSlug string
	var topLinkClicks int64
	err = r.db.QueryRow(context.Background(),
		topLinkQuery,
		userID,
	).Scan(
		&topLinkTitle,
		&topLinkSlug,
		&topLinkClicks,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			overview.TopLink = nil
		} else {
			logger.Log.Error(
				"failed to query top link for overview",
				"error", err,
			)
			return nil, err
		}
	} else {
		overview.TopLink = &models.TopLinkDTO{
			Title:  topLinkTitle,
			Slug:   topLinkSlug,
			Clicks: topLinkClicks,
		}
	}

	return &overview, nil
}

func (r *linkRepository) GetLinkAnalytics(userID uuid.UUID) ([]models.LinkAnalyticsDTO, error) {
	query := `
	SELECT title, slug, click_count 
	FROM links WHERE user_id = $1 ORDER BY click_count DESC
	`
	rows, err := r.db.Query(
		context.Background(),
		query,
		userID,
	)
	if err != nil {
		logger.Log.Error(
			"failed to query link analytics",
			"error", err,
		)
		return nil, err
	}
	defer rows.Close()
	analytics := []models.LinkAnalyticsDTO{}
	for rows.Next() {
		var item models.LinkAnalyticsDTO
		err := rows.Scan(
			&item.Title,
			&item.Slug,
			&item.Clicks,
		)
		if err != nil {
			logger.Log.Error(
				"failed to scan link analytics row",
				"error", err,
			)
			return nil, err
		}
		analytics = append(analytics, item)
	}
	if err = rows.Err(); err != nil {
		logger.Log.Error(
			"failed during link analytics rows iteration",
			"error", err,
		)
		return nil, err
	}
	return analytics, nil
}

func (r *linkRepository) GetDailyAnalytics(userID uuid.UUID, days int) ([]models.DailyAnalyticsDTO, error) {
	query := `
	SELECT
		TO_CHAR(ce.clicked_at, 'YYYY-MM-DD') as date,
		COUNT(*) as clicks
	FROM click_events ce
	JOIN links l ON ce.link_id = l.id
	WHERE l.user_id = $1 AND ce.clicked_at >= CURRENT_DATE - ($2 * INTERVAL '1 day')
	GROUP BY TO_CHAR(ce.clicked_at, 'YYYY-MM-DD')
	ORDER BY TO_CHAR(ce.clicked_at, 'YYYY-MM-DD') DESC;
	`
	rows, err := r.db.Query(
		context.Background(),
		query,
		userID,
		days,
	)
	if err != nil {
		logger.Log.Error(
			"failed to query daily analytics",
			"error", err,
		)
		return nil, err
	}
	defer rows.Close()
	analytics := []models.DailyAnalyticsDTO{}
	for rows.Next() {
		var item models.DailyAnalyticsDTO
		err := rows.Scan(
			&item.Date,
			&item.Clicks,
		)
		if err != nil {
			logger.Log.Error(
				"failed to scan daily analytics row",
				"error", err,
			)
			return nil, err
		}
		analytics = append(analytics, item)
	}
	if err = rows.Err(); err != nil {
		logger.Log.Error(
			"failed during daily analytics rows iteration",
			"error", err,
		)
		return nil, err
	}
	return analytics, nil
}

func (r *linkRepository) GetRecentActivity(userID uuid.UUID, limit int) ([]models.RecentActivityDTO, error) {
	query := `
	SELECT l.title, l.slug, ce.clicked_at, ce.ip_address, ce.referrer
	FROM click_events ce
	JOIN links l ON ce.link_id = l.id
	WHERE l.user_id = $1
	ORDER BY ce.clicked_at DESC
	LIMIT $2;
	`
	rows, err := r.db.Query(
		context.Background(),
		query,
		userID,
		limit,
	)
	if err != nil {
		logger.Log.Error(
			"failed to query recent activity",
			"error", err,
		)
		return nil, err
	}
	defer rows.Close()
	activities := []models.RecentActivityDTO{}
	for rows.Next() {
		var item models.RecentActivityDTO
		err := rows.Scan(
			&item.LinkTitle,
			&item.Slug,
			&item.ClickedAt,
			&item.IPAddress,
			&item.Referrer,
		)
		if err != nil {
			logger.Log.Error(
				"failed to scan recent activity row",
				"error", err,
			)
			return nil, err
		}
		activities = append(activities, item)
	}
	if err = rows.Err(); err != nil {
		logger.Log.Error(
			"failed during recent activity rows iteration",
			"error", err,
		)
		return nil, err
	}
	return activities, nil
}

func (r *linkRepository) GetReferrerAnalytics(
	userID uuid.UUID,
) ([]models.ReferrerAnalyticsDTO, error) {
	query := `
		SELECT COALESCE(NULLIF(ce.referrer, ''), 'Direct') AS source,
		COUNT(*) AS clicks
		FROM click_events ce
		JOIN links l ON ce.link_id = l.id
		WHERE l.user_id = $1
		GROUP BY source
		ORDER BY clicks DESC;
	`
	rows, err := r.db.Query(
		context.Background(),
		query,
		userID,
	)
	if err != nil {
		logger.Log.Error(
			"failed to query referrer analytics",
			"error", err,
		)
		return nil, err
	}
	defer rows.Close()
	aggregated := make(map[string]int64)
	for rows.Next() {
		var item models.ReferrerAnalyticsDTO
		err := rows.Scan(
			&item.Source,
			&item.Clicks,
		)
		if err != nil {
			logger.Log.Error(
				"failed to scan referrer analytics row",
				"error", err,
			)
			return nil, err
		}
		normSource := utils.NormalizeReferrer(item.Source)
		aggregated[normSource] += item.Clicks
	}
	if err = rows.Err(); err != nil {
		logger.Log.Error(
			"failed during referrer analytics rows iteration",
			"error", err,
		)
		return nil, err
	}

	analytics := []models.ReferrerAnalyticsDTO{}
	for source, clicks := range aggregated {
		analytics = append(analytics, models.ReferrerAnalyticsDTO{
			Source: source,
			Clicks: clicks,
		})
	}

	sort.Slice(analytics, func(i, j int) bool {
		return analytics[i].Clicks > analytics[j].Clicks
	})

	return analytics, nil
}

