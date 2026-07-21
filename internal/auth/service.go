package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/MRNaveed-stack/LinkPulse/internal/logger"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/repository"
	"github.com/MRNaveed-stack/LinkPulse/internal/utils"
)

var (
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidEmail    = errors.New("invalid email or password")
	ErrInvalidPassword = errors.New("invalid email or password")
	ErrUserNotFound    = errors.New("user not found")
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.LoginResponse, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token string, newPassword string) error
	ValidateToken(token string) (*TokenClaims, error)
	LoginOrRegisterWithGoogle(ctx context.Context, googleID, email, name string) (*models.LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*models.LoginResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

type authService struct {
	userRepo               repository.UserRepository
	passwordResetTokenRepo repository.PasswordResetTokenRepository
	profileRepo            repository.ProfileRepository
}

func NewAuthService(
	userRepo repository.UserRepository,
	prRepo repository.PasswordResetTokenRepository,
	profileRepo repository.ProfileRepository,
) AuthService {
	return &authService{
		userRepo:               userRepo,
		passwordResetTokenRepo: prRepo,
		profileRepo:            profileRepo,
	}
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) (*models.LoginResponse, error) {
	// Check if user exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate tokens first (safely outside DB operations)
	userID := uuid.New()
	accessToken, refreshToken, err := GenerateTokenPair(userID, req.Email, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	// Create user
	user := &models.User{
		ID:           userID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Plan:         "free",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Log.Info(
		"user registered",
		"user_id", user.ID,
		"email", user.Email,
	)

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Token:        accessToken,
	}, nil
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidEmail
	}

	// Check password
	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidPassword
	}

	logger.Log.Info(
		"user logged in",
		"user_id", user.ID,
	)

	accessToken, refreshToken, err := GenerateTokenPair(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Token:        accessToken,
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	opaqueToken, err := generateResetToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Persist token server-side
	expHours := os.Getenv("PASSWORD_RESET_EXPIRATION_HOURS")
	if expHours == "" {
		expHours = "1" // default 1 hour
	}
	expiration, err := time.ParseDuration(expHours + "h")
	if err != nil {
		return "", fmt.Errorf("invalid PASSWORD_RESET_EXPIRATION_HOURS: %w", err)
	}

	model := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     opaqueToken,
		ExpiresAt: time.Now().Add(expiration),
		Used:      false,
		CreatedAt: time.Now(),
	}

	if err := s.passwordResetTokenRepo.Create(ctx, model); err != nil {
		return "", fmt.Errorf("failed to create reset token: %w", err)
	}

	return opaqueToken, nil
}

func (s *authService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	resetToken, err := s.passwordResetTokenRepo.GetByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid reset token: %w", err)
	}

	if resetToken.Used {
		return fmt.Errorf("reset token already used")
	}
	if time.Now().After(resetToken.ExpiresAt) {
		return fmt.Errorf("reset token expired")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, resetToken.UserID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := s.passwordResetTokenRepo.MarkUsed(ctx, resetToken.ID, time.Now()); err != nil {
		return fmt.Errorf("failed to mark reset token used: %w", err)
	}

	return nil
}

func (s *authService) ValidateToken(token string) (*TokenClaims, error) {
	claims, err := ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// generateResetToken creates a cryptographically secure random token
func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// TokenClaims is an alias for auth.Claims for consistency
type TokenClaims = Claims

func (s *authService) LoginOrRegisterWithGoogle(ctx context.Context, googleID, email, name string) (*models.LoginResponse, error) {
	// 1. Check if user exists by Google ID
	user, err := s.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user by Google ID: %w", err)
	}

	if user != nil {
		// User exists, login
		accessToken, refreshToken, err := GenerateTokenPair(user.ID, user.Email, user.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to generate token pair: %w", err)
		}
		logger.Log.Info(
			"user logged in",
			"user_id", user.ID,
		)
		return &models.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Token:        accessToken,
		}, nil
	}

	// 2. Check if user exists by Email
	user, err = s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user by email: %w", err)
	}

	if user != nil {
		// Link Google ID to existing account
		err = s.userRepo.LinkGoogleAccount(ctx, user.ID, googleID)
		if err != nil {
			return nil, fmt.Errorf("failed to link Google account: %w", err)
		}

		accessToken, refreshToken, err := GenerateTokenPair(user.ID, user.Email, user.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to generate token pair: %w", err)
		}
		logger.Log.Info(
			"user logged in",
			"user_id", user.ID,
		)
		return &models.LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Token:        accessToken,
		}, nil
	}

	// 3. User does not exist, register new user
	// Generate unique username
	username, err := s.generateUniqueUsername(ctx, email, name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate unique username: %w", err)
	}

	// Generate a cryptographically secure random password hash (inaccessible by standard login)
	tempPassBytes := make([]byte, 32)
	if _, err := rand.Read(tempPassBytes); err != nil {
		return nil, fmt.Errorf("failed to generate temporary password: %w", err)
	}
	hashedPassword, err := utils.HashPassword(hex.EncodeToString(tempPassBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &models.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Plan:         "free",
		GoogleID:     &googleID,
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Log.Info(
		"user registered",
		"user_id", newUser.ID,
		"email", newUser.Email,
	)

	// Create a profile for the user
	displayName := name
	if displayName == "" {
		displayName = username
	}
	profile := &models.Profile{
		UserID:      newUser.ID,
		DisplayName: displayName,
		Bio:         "Welcome to my LinkPulse profile!",
		AvatarURL:   "",
	}

	if err := s.profileRepo.CreateProfile(profile); err != nil {
		return nil, fmt.Errorf("failed to create default user profile: %w", err)
	}

	// Generate token for the new user
	accessToken, refreshToken, err := GenerateTokenPair(newUser.ID, newUser.Email, newUser.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Token:        accessToken,
	}, nil
}

func (s *authService) generateUniqueUsername(ctx context.Context, email, fullName string) (string, error) {
	// First choice: email prefix
	parts := strings.Split(email, "@")
	base := parts[0]
	if base == "" && fullName != "" {
		base = fullName
	}

	// Sanitize base: allow only alphanumeric
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	base = reg.ReplaceAllString(base, "")
	base = strings.ToLower(base)

	if len(base) < 2 {
		base = "user"
	}
	if len(base) > 30 {
		base = base[:30]
	}

	username := base
	// Query up to 10 candidates with suffix
	for i := 0; i < 10; i++ {
		existing, err := s.userRepo.GetByUsername(ctx, username)
		if err != nil {
			return "", err
		}
		if existing == nil {
			return username, nil
		}
		// Append random suffix
		username = fmt.Sprintf("%s%d", base, time.Now().UnixNano()%9000+1000)
	}

	return fmt.Sprintf("%s%d", base, time.Now().UnixNano()%100000), nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*models.LoginResponse, error) {
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in claims: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	accessToken, newRefreshToken, err := GenerateTokenPair(userID, claims.Email, claims.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		Token:        accessToken,
	}, nil
}

func (s *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
