package service

import (
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/repository"
	"github.com/google/uuid"
)

type AnalyticsService interface {
	GetOverview(userID uuid.UUID) (*models.AnalyticsOverview, error)
	GetLinkAnalytics(userID uuid.UUID) ([]models.LinkAnalyticsDTO, error)
	GetDailyAnalytics(userID uuid.UUID, days int) ([]models.DailyAnalyticsDTO, error)
	GetRecentActivity(userID uuid.UUID, limit int) ([]models.RecentActivityDTO, error)
	GetReferrerAnalytics(userID uuid.UUID) ([]models.ReferrerAnalyticsDTO, error)
}

type analyticsService struct {
	linkRepo repository.LinkRepository
}

func NewAnalyticsService(linkRepo repository.LinkRepository) AnalyticsService {
	return &analyticsService{
		linkRepo: linkRepo,
	}
}

func (s *analyticsService) GetOverview(
	userID uuid.UUID,
) (*models.AnalyticsOverview, error) {
	return s.linkRepo.GetAnalyticsOverview(
		userID,
	)
}

func (s *analyticsService) GetLinkAnalytics(
	userID uuid.UUID,
) ([]models.LinkAnalyticsDTO, error) {
	return s.linkRepo.GetLinkAnalytics(
		userID,
	)
}
func (s *analyticsService) GetDailyAnalytics(
	userID uuid.UUID,
	days int,
) ([]models.DailyAnalyticsDTO, error) {
	return s.linkRepo.GetDailyAnalytics(
		userID,
		days,
	)

}

func (s *analyticsService) GetRecentActivity(
	userID uuid.UUID,
	limit int,
) ([]models.RecentActivityDTO, error) {

	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Maximum limit
	}
	return s.linkRepo.GetRecentActivity(userID, limit)

}

func (s *analyticsService) GetReferrerAnalytics(
	userID uuid.UUID,
) ([]models.ReferrerAnalyticsDTO, error) {
	return s.linkRepo.GetReferrerAnalytics(userID)
}
