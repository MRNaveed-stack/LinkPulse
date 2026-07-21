package service_test

import (
	"context"
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProfileService_GetPublicProfile(t *testing.T) {
	username := "testuser"
	userID := uuid.New()

	userRepo := &MockUserRepository{
		GetByUsernameFunc: func(ctx context.Context, uname string) (*models.User, error) {
			assert.Equal(t, username, uname)
			return &models.User{ID: userID, Username: username}, nil
		},
	}

	profileRepo := &MockProfileRepository{
		GetByUserIDFunc: func(id uuid.UUID) (*models.Profile, error) {
			assert.Equal(t, userID, id)
			return &models.Profile{
				UserID:      userID,
				DisplayName: "Test User Profile",
				Bio:         "Golang engineer",
			}, nil
		},
	}

	linkRepo := &MockLinkRepository{
		GetActiveLinksByUserIDFunc: func(id uuid.UUID) ([]*models.Link, error) {
			assert.Equal(t, userID, id)
			return []*models.Link{
				{Title: "Github", Slug: "gitslug"},
			}, nil
		},
	}

	svc := service.NewProfileService(userRepo, profileRepo, linkRepo)
	profileResp, err := svc.GetPublicProfile(username)
	assert.NoError(t, err)
	assert.Equal(t, "Test User Profile", profileResp.DisplayName)
	assert.Len(t, profileResp.Links, 1)
}

func TestProfileService_UpdateProfile(t *testing.T) {
	userID := uuid.New()
	req := &models.UpdateProfileRequest{
		DisplayName: "New Display Name",
		Bio:         "New Bio",
		AvatarURL:   "http://avatar.url",
	}

	displayName := "Old Display Name"
	profileRepo := &MockProfileRepository{
		GetByUserIDFunc: func(id uuid.UUID) (*models.Profile, error) {
			// Profile exists
			return &models.Profile{UserID: userID, DisplayName: displayName}, nil
		},
		UpdateProfileFunc: func(p *models.Profile) error {
			assert.Equal(t, req.DisplayName, p.DisplayName)
			assert.Equal(t, req.Bio, p.Bio)
			displayName = p.DisplayName
			return nil
		},
	}
	userRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.User, error) {
			assert.Equal(t, userID, id)
			return &models.User{
				ID:       userID,
				Username: "testuser",
			}, nil
		},
	}

	svc := service.NewProfileService(userRepo, profileRepo, nil)
	updated, err := svc.UpdateProfile(userID, req)
	assert.NoError(t, err)
	assert.Equal(t, req.DisplayName, updated.DisplayName)
}
