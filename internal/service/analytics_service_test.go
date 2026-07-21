package service_test

import (
	"testing"
	"time"

	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAndAnalyticsService_GetOverview(t *testing.T) {
	userID := uuid.New()
	Overview := &models.AnalyticsOverview{
		TotalLinks:  5,
		TotalClicks: 150,
		ClicksToday: 10,
		ActiveLinks: 4,
		TopLink: &models.TopLinkDTO{
			Title:  "Google link",
			Slug:   "google",
			Clicks: 7,
		},
	}
	linkRepo := &MockLinkRepository{
		GetAnalyticsOverviewFunc: func(id uuid.UUID) (*models.AnalyticsOverview,
			error) {
			assert.Equal(t, userID, id)
			return Overview, nil
		},
	}
	svc := service.NewAnalyticsService(linkRepo)
	res, err := svc.GetOverview(userID)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), res.TotalLinks)

}

func TestAnalyticsService_GetLinkAnalytics(t *testing.T) {
	userID := uuid.New()
	expected := []models.LinkAnalyticsDTO{
		{Title: "Title", Slug: "slug", Clicks: 50},
	}
	linkRepo := &MockLinkRepository{
		GetLinkAnalyticsFunc: func(id uuid.UUID) ([]models.LinkAnalyticsDTO,
			error) {
			assert.Equal(t, userID, id)
			return expected, nil
		},
	}
	svc := service.NewAnalyticsService(linkRepo)
	res, err := svc.GetLinkAnalytics(userID)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestAnalyticsService_GetDailyAnalytics(t *testing.T) {
	userID := uuid.New()
	expected := []models.DailyAnalyticsDTO{
		{Date: "2026-07-01", Clicks: 50},
	}
	linkRepo := &MockLinkRepository{
		GetDailyAnalyticsFunc: func(id uuid.UUID, days int) ([]models.DailyAnalyticsDTO, error) {
			assert.Equal(t, 7, days)
			return expected, nil
		},
	}
	svc := service.NewAnalyticsService(linkRepo)
	res, err := svc.GetDailyAnalytics(userID, 7)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestAnalyticsService_GetRecentActivity(t *testing.T) {
	userID := uuid.New()
	expected := []models.RecentActivityDTO{
		{LinkTitle: "Google", Slug: "google", ClickedAt: time.Now()},
	}
	linkRepo := &MockLinkRepository{
		GetRecentActivityFunc: func(id uuid.UUID, limit int) ([]models.RecentActivityDTO, error) {
			assert.Equal(t, 20, limit)
			return expected, nil
		},
	}
	svc := service.NewAnalyticsService(linkRepo)
	res, err := svc.GetRecentActivity(userID, 0)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestAnalyticsService_GetReferrerAnalytics(t *testing.T) {
	userID := uuid.New()
	expected := []models.ReferrerAnalyticsDTO{
		{Source: "Google", Clicks: 20},
	}
	linkRepo := &MockLinkRepository{
		GetReferrerAnalyticsFunc: func(id uuid.UUID) ([]models.ReferrerAnalyticsDTO, error) {
			return expected, nil
		},
	}
	svc := service.NewAnalyticsService(linkRepo)
	res, err := svc.GetReferrerAnalytics(userID)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
}
