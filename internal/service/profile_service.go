package service

import (
	"context"
	"fmt"

	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/repository"
	"github.com/MRNaveed-stack/LinkPulse/internal/utils"
	"github.com/google/uuid"
)

type ProfileService interface {
	GetPublicProfile(username string) (*models.PublicProfileResponse, error)
	UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) (*models.Profile, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error
	DeleteAccount(ctx context.Context, userID uuid.UUID) error
}

type profileService struct {
	userRepo    repository.UserRepository
	profileRepo repository.ProfileRepository
	linkRepo    repository.LinkRepository
}

func NewProfileService(
	userRepo repository.UserRepository,
	profileRepo repository.ProfileRepository,
	linkRepo repository.LinkRepository,
) ProfileService {
	return &profileService{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		linkRepo:    linkRepo,
	}
}

func (s *profileService) GetPublicProfile(username string) (*models.PublicProfileResponse, error) {
	user, err := s.userRepo.GetByUsername(context.Background(), username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user %q not found", username)
	}
	profile, err := s.profileRepo.GetByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	links, err := s.linkRepo.GetActiveLinksByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	var displayName, bio, avatarURL string
	if profile != nil {
		displayName = profile.DisplayName
		bio = profile.Bio
		avatarURL = profile.AvatarURL
	} else {
		displayName = user.Username
	}

	response := &models.PublicProfileResponse{
		Username:    user.Username,
		DisplayName: displayName,
		Bio:         bio,
		AvatarURL:   avatarURL,
	}

	for _, link := range links {
		response.Links = append(
			response.Links,
			models.PublicLinkDTO{
				Title: link.Title,
				Slug:  link.Slug,
			},
		)
	}
	return response, nil
}

func (s *profileService) UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) (*models.Profile, error) {
	ctx := context.Background()

	// Update user details if different
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if req.Username != "" && req.Username != user.Username {
		existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username: %w", err)
		}
		if existingUser != nil {
			return nil, fmt.Errorf("username %q is already taken", req.Username)
		}
		err = s.userRepo.UpdateUsername(ctx, userID, req.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to update username: %w", err)
		}
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if existingUser != nil {
			return nil, fmt.Errorf("email %q is already taken", req.Email)
		}
		err = s.userRepo.UpdateEmail(ctx, userID, req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to update email: %w", err)
		}
	}

	existingProfile, err := s.profileRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	var profile *models.Profile
	if existingProfile == nil {
		profile = &models.Profile{
			UserID:      userID,
			DisplayName: req.DisplayName,
			Bio:         req.Bio,
			AvatarURL:   req.AvatarURL,
		}
		err = s.profileRepo.CreateProfile(profile)
		if err != nil {
			return nil, fmt.Errorf("failed to create profile: %w", err)
		}
	} else {
		existingProfile.DisplayName = req.DisplayName
		existingProfile.Bio = req.Bio
		existingProfile.AvatarURL = req.AvatarURL
		err = s.profileRepo.UpdateProfile(existingProfile)
		if err != nil {
			return nil, fmt.Errorf("failed to update profile: %w", err)
		}
		profile = existingProfile
	}
	return s.profileRepo.GetByUserID(userID)
}

func (s *profileService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	if err := utils.CheckPassword(currentPassword, user.PasswordHash); err != nil {
		return fmt.Errorf("invalid current password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Save to DB
	return s.userRepo.UpdatePassword(ctx, userID, hashedPassword)
}

func (s *profileService) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(ctx, userID)
}

