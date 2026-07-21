package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/handler"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAnalyticsHandler_GetOverview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAnalyticsService{}
	h := handler.NewAnalyticsHandler(mockSvc)

	userID := uuid.New()
	mockSvc.GetOverviewFunc = func(uid uuid.UUID) (*models.AnalyticsOverview, error) {
		return &models.AnalyticsOverview{TotalLinks: 10}, nil
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/analytics/overview", nil)

	_, r := gin.CreateTestContext(w)
	r.GET("/analytics/overview", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.GetOverview)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnalyticsHandler_GetLinkAnalytics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAnalyticsService{}
	h := handler.NewAnalyticsHandler(mockSvc)

	userID := uuid.New()
	mockSvc.GetLinkAnalyticsFunc = func(uid uuid.UUID) ([]models.LinkAnalyticsDTO, error) {
		return []models.LinkAnalyticsDTO{{Title: "L1"}}, nil
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/analytics/links", nil)

	_, r := gin.CreateTestContext(w)
	r.GET("/analytics/links", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.GetLinkAnalytics)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnalyticsHandler_GetDailyAnalytics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAnalyticsService{}
	h := handler.NewAnalyticsHandler(mockSvc)

	userID := uuid.New()
	mockSvc.GetDailyAnalyticsFunc = func(uid uuid.UUID, days int) ([]models.DailyAnalyticsDTO, error) {
		assert.Equal(t, 7, days)
		return []models.DailyAnalyticsDTO{{Clicks: 10}}, nil
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/analytics/daily?days=7", nil)

	_, r := gin.CreateTestContext(w)
	r.GET("/analytics/daily", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.GetDailyAnalytics)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnalyticsHandler_GetRecentActivity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAnalyticsService{}
	h := handler.NewAnalyticsHandler(mockSvc)

	userID := uuid.New()
	mockSvc.GetRecentActivityFunc = func(uid uuid.UUID, limit int) ([]models.RecentActivityDTO, error) {
		assert.Equal(t, 10, limit)
		return []models.RecentActivityDTO{{Slug: "l1"}}, nil
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/analytics/activity?limit=10", nil)

	_, r := gin.CreateTestContext(w)
	r.GET("/analytics/activity", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.GetRecentActivity)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAnalyticsHandler_GetReferrerAnalytics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAnalyticsService{}
	h := handler.NewAnalyticsHandler(mockSvc)

	userID := uuid.New()
	mockSvc.GetReferrerAnalyticsFunc = func(uid uuid.UUID) ([]models.ReferrerAnalyticsDTO, error) {
		return []models.ReferrerAnalyticsDTO{{Source: "Direct"}}, nil
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/analytics/referrers", nil)

	_, r := gin.CreateTestContext(w)
	r.GET("/analytics/referrers", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.GetReferrerAnalytics)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
