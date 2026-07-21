package handler

import (
	"net/http"

	"github.com/MRNaveed-stack/LinkPulse/internal/logger"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsService interface {
	GetOverview(userID uuid.UUID) (*models.AnalyticsOverview, error)
	GetLinkAnalytics(userID uuid.UUID) ([]models.LinkAnalyticsDTO, error)
	GetDailyAnalytics(userID uuid.UUID, days int) ([]models.DailyAnalyticsDTO, error)
	GetRecentActivity(userID uuid.UUID, limit int) ([]models.RecentActivityDTO, error)
	GetReferrerAnalytics(userID uuid.UUID) ([]models.ReferrerAnalyticsDTO, error)
}

type AnalyticsHandler struct {
	service AnalyticsService
}

func NewAnalyticsHandler(service AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}

func (h *AnalyticsHandler) GetOverview(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	overview, err := h.service.GetOverview(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get analytics overview"})
		return
	}
	c.JSON(http.StatusOK, overview)
}

func (h *AnalyticsHandler) GetLinkAnalytics(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	data, err := h.service.GetLinkAnalytics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *AnalyticsHandler) GetDailyAnalytics(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	var query models.AnalyticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	analytics, err := h.service.GetDailyAnalytics(userID, query.Days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get daily analytics"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": analytics})
}

func (h *AnalyticsHandler) GetRecentActivity(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	var query models.RecentActivityQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	activity, err := h.service.GetRecentActivity(userID, query.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recent activity"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": activity})
}

func (h *AnalyticsHandler) GetReferrerAnalytics(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)
	data, err := h.service.GetReferrerAnalytics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
