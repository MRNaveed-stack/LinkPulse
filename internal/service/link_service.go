package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/repository"
	"github.com/google/uuid"
)

type LinkService interface {
	CreateLink(
		userID uuid.UUID,
		req models.CreateLinkRequest,
	) error
	GetUserLinks(
		userID uuid.UUID,
	) ([]*models.Link, error)
	HandleRedirect(
		userID uuid.UUID,
		slug string,
		ip string,
		userAgent string,
		referrer string,
	) (string, error)
	UpdateLink(
		userID uuid.UUID,
		linkID uuid.UUID,
		req models.UpdateLinkRequest,
	) (*models.Link, error)
	DeleteLink(
		userID uuid.UUID,
		linkID uuid.UUID,
	) error
	UpdateLinkStatus(
		userID uuid.UUID,
		linkID uuid.UUID,
		isActive bool,
	) error
}

type linkService struct {
	linkRepo  repository.LinkRepository
	clickRepo repository.ClickRepository
}

func NewLinkService(
	linkRepo repository.LinkRepository,
	clickRepo repository.ClickRepository,
) LinkService {
	return &linkService{
		linkRepo:  linkRepo,
		clickRepo: clickRepo,
	}
}

func (s *linkService) CreateLink(userID uuid.UUID, req models.CreateLinkRequest) error {
	link := &models.Link{
		ID:             uuid.New(),
		UserID:         userID,
		Title:          req.Title,
		Slug:           req.Slug,
		DestinationURL: req.DestinationURL,
		IsActive:       true,
		ClickCount:     0,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return s.linkRepo.Create(link)
}

func (s *linkService) GetUserLinks(
	userID uuid.UUID,
) ([]*models.Link, error) {
	links, err := s.linkRepo.GetByUserID(userID)
	return links, err
}

func (s *linkService) HandleRedirect(
	userID uuid.UUID,
	slug string,
	ip string,
	userAgent string,
	referrer string,
) (string, error) {
	link, err := s.linkRepo.GetByUserIDAndSlug(userID, slug)
	if err != nil {
		return "", err
	}
	if link == nil {
		return "", errors.New("link not found")
	}
	if !link.IsActive {
		return "", errors.New("link is disabled")
	}
	event := &models.ClickEvent{
		ID:        uuid.New(),
		LinkID:    link.ID,
		IPAddress: ip,
		UserAgent: userAgent,
		Referrer:  referrer,
		ClickedAt: time.Now(),
	}

	if err := s.clickRepo.Create(event); err != nil {
		return "", err
	}
	if err := s.linkRepo.IncrementClickCount(link.ID); err != nil {
		return "", err
	}
	return link.DestinationURL, nil
}

func (s *linkService) UpdateLink(
	userID uuid.UUID,
	linkID uuid.UUID,
	req models.UpdateLinkRequest,
) (*models.Link, error) {
	link, err := s.linkRepo.GetByID(linkID)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, errors.New("link not found")
	}
	if link.UserID != userID {
		return nil, errors.New("unauthorized: you don't own this link")
	}
	link.Title = req.Title
	link.Slug = req.Slug
	link.DestinationURL = req.DestinationURL

	err = s.linkRepo.Update(link)
	if err != nil {
		return nil, fmt.Errorf("failed to update link: %w", err)
	}
	return link, nil
}

func (s *linkService) DeleteLink(userID uuid.UUID, linkID uuid.UUID) error {
	link, err := s.linkRepo.GetByID(linkID)
	if err != nil {
		return fmt.Errorf("failed to retrieve link: %w", err)
	}
	if link == nil {
		return errors.New("link not found")
	}
	if link.UserID != userID {
		return errors.New("unauthorized: you don't own this link")
	}
	err = s.linkRepo.Delete(linkID)
	if err != nil {
		return fmt.Errorf("failed to delete link: %w", err)
	}
	return nil
}

func (s *linkService) UpdateLinkStatus(userID uuid.UUID, linkID uuid.UUID, isActive bool) error {
	link, err := s.linkRepo.GetByID(linkID)
	if err != nil {
		return fmt.Errorf("link not found: %w", err)
	}
	if link == nil {
		return errors.New("link not found")
	}
	if link.UserID != userID {
		return errors.New("unauthorized: you don't own this link")
	}
	err = s.linkRepo.UpdateStatus(linkID, isActive)
	if err != nil {
		return fmt.Errorf("failed to update link status: %w", err)
	}
	return nil
}
