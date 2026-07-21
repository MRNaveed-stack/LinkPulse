package service_test

import (
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLinkService_CreateLink(t *testing.T) {
	userID := uuid.New()
	req := models.CreateLinkRequest{
		Title:          "My Link",
		Slug:           "mylink",
		DestinationURL: "https://google.com",
	}
	linkRepo := &MockLinkRepository{
		CreateFunc: func(link *models.Link) error {
			assert.Equal(t, userID, link.UserID)
			assert.Equal(t, req.Title, link.Title)
			assert.Equal(t, req.Slug, link.Slug)
			return nil
		},
	}
	svc := service.NewLinkService(linkRepo, nil)
	err := svc.CreateLink(userID, req)
	assert.NoError(t, err)
}
func TestLinkService_GetUserLinks(t *testing.T) {
	userID := uuid.New()
	expectedLinks := []*models.Link{
		{ID: uuid.New(), UserID: userID, Title: "Link 1"},
		{ID: uuid.New(), UserID: userID, Title: "Link 2"},
	}
	linkRepo := &MockLinkRepository{
		GetByUserIDFunc: func(id uuid.UUID) ([]*models.Link, error) {
			assert.Equal(t, userID, id)
			return expectedLinks, nil
		},
	}
	svc := service.NewLinkService(linkRepo, nil)
	links, err := svc.GetUserLinks(userID)
	assert.NoError(t, err)
	assert.Len(t, links, 2)
}
func TestLinkService_HandleRedirect(t *testing.T) {
	userID := uuid.New()
	slug := "redirectslug"
	destination := "https://destination.com"
	linkID := uuid.New()
	linkRepo := &MockLinkRepository{
		GetByUserIDAndSlugFunc: func(uid uuid.UUID, s string) (*models.Link, error) {
			assert.Equal(t, userID, uid)
			assert.Equal(t, slug, s)
			return &models.Link{
				ID:             linkID,
				UserID:         userID,
				Slug:           slug,
				DestinationURL: destination,
				IsActive:       true,
			}, nil
		},
		IncrementClickCountFunc: func(id uuid.UUID) error {
			assert.Equal(t, linkID, id)
			return nil
		},
	}
	clickRepo := &MockClickRepository{
		CreateFunc: func(event *models.ClickEvent) error {
			assert.Equal(t, linkID, event.LinkID)
			assert.Equal(t, "192.168.1.1", event.IPAddress)
			return nil
		},
	}
	svc := service.NewLinkService(linkRepo, clickRepo)
	dest, err := svc.HandleRedirect(userID, slug, "192.168.1.1", "Mozilla/5.0",
		"https://referrer.com")
	assert.NoError(t, err)
	assert.Equal(t, destination, dest)
}
func TestLinkService_UpdateLink(t *testing.T) {
	userID := uuid.New()
	linkID := uuid.New()
	req := models.UpdateLinkRequest{
		Title:          "Updated Title",
		Slug:           "updatedslug",
		DestinationURL: "https://newurl.com",
	}
	linkRepo := &MockLinkRepository{
		GetByIDFunc: func(id uuid.UUID) (*models.Link, error) {
			return &models.Link{ID: linkID, UserID: userID, Title: "Old Title"},
				nil
		},
		UpdateFunc: func(link *models.Link) error {
			assert.Equal(t, req.Title, link.Title)
			return nil
		},
	}
	svc := service.NewLinkService(linkRepo, nil)
	updatedLink, err := svc.UpdateLink(userID, linkID, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Title, updatedLink.Title)
}
func TestLinkService_DeleteLink(t *testing.T) {
	userID := uuid.New()
	linkID := uuid.New()
	linkRepo := &MockLinkRepository{
		GetByIDFunc: func(id uuid.UUID) (*models.Link, error) {
			return &models.Link{ID: linkID, UserID: userID}, nil
		},
		DeleteFunc: func(id uuid.UUID) error {
			assert.Equal(t, linkID, id)
			return nil
		},
	}
	svc := service.NewLinkService(linkRepo, nil)
	err := svc.DeleteLink(userID, linkID)
	assert.NoError(t, err)
}
func TestLinkService_UpdateLinkStatus(t *testing.T) {
	userID := uuid.New()
	linkID := uuid.New()
	linkRepo := &MockLinkRepository{
		GetByIDFunc: func(id uuid.UUID) (*models.Link, error) {
			return &models.Link{ID: linkID, UserID: userID}, nil
		},
		UpdateStatusFunc: func(id uuid.UUID, active bool) error {
			assert.Equal(t, linkID, id)
			assert.True(t, active)
			return nil
		},
	}
	svc := service.NewLinkService(linkRepo, nil)
	err := svc.UpdateLinkStatus(userID, linkID, true)
	assert.NoError(t, err)
}
