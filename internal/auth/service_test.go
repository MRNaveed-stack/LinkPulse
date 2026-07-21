package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/MRNaveed-stack/LinkPulse/internal/auth"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Setup secrets for the test suite
func init() {
	_ = auth.Claims{} // Verify imports compiled
}

func TestAuthService_Register(t *testing.T) {
	t.Setenv("JWT_SECRET", "mysecret")
	t.Setenv("JWT_REFRESH_SECRET", "myrefreshsecret")

	ctx := context.Background()
	req := models.RegisterRequest{
		Username: "newuser",
		Email:    "new@linkpulse.com",
		Password: "password123",
	}

	// Success case
	userRepo := &MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return nil, nil // User does not exist
		},
		CreateFunc: func(ctx context.Context, user *models.User) error {
			assert.Equal(t, req.Username, user.Username)
			assert.Equal(t, req.Email, user.Email)
			return nil
		},
	}
	prRepo := &MockPasswordResetTokenRepository{}
	profileRepo := &MockProfileRepository{}

	svc := auth.NewAuthService(userRepo, prRepo, profileRepo)
	resp, err := svc.Register(ctx, req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)

	// User exists case
	userRepo.GetByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		return &models.User{ID: uuid.New(), Email: email}, nil
	}
	_, err = svc.Register(ctx, req)
	assert.ErrorIs(t, err, auth.ErrUserExists)
}

func TestAuthService_Login(t *testing.T) {
	t.Setenv("JWT_SECRET", "mysecret")
	t.Setenv("JWT_REFRESH_SECRET", "myrefreshsecret")

	ctx := context.Background()
	req := models.LoginRequest{
		Email:    "test@linkpulse.com",
		Password: "password123",
	}

	hashedPassword, _ := bcryptHash("password123")
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	// Successful login
	userRepo := &MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return user, nil
		},
	}
	svc := auth.NewAuthService(userRepo, nil, nil)
	resp, err := svc.Login(ctx, req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)

	// Incorrect password
	req.Password = "wrongpassword"
	_, err = svc.Login(ctx, req)
	assert.ErrorIs(t, err, auth.ErrInvalidPassword)

	// User not found
	userRepo.GetByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		return nil, nil
	}
	_, err = svc.Login(ctx, req)
	assert.ErrorIs(t, err, auth.ErrInvalidEmail)
}

func TestAuthService_ForgotPassword(t *testing.T) {
	ctx := context.Background()
	email := "user@linkpulse.com"
	userID := uuid.New()

	userRepo := &MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return &models.User{ID: userID, Email: email}, nil
		},
	}
	prRepo := &MockPasswordResetTokenRepository{
		CreateFunc: func(ctx context.Context, token *models.PasswordResetToken) error {
			assert.Equal(t, userID, token.UserID)
			assert.NotEmpty(t, token.Token)
			return nil
		},
	}

	svc := auth.NewAuthService(userRepo, prRepo, nil)
	token, err := svc.ForgotPassword(ctx, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// User not found
	userRepo.GetByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		return nil, nil
	}
	_, err = svc.ForgotPassword(ctx, email)
	assert.ErrorIs(t, err, auth.ErrUserNotFound)
}

func TestAuthService_ResetPassword(t *testing.T) {
	ctx := context.Background()
	token := "reset-token"
	userID := uuid.New()

	prRepo := &MockPasswordResetTokenRepository{
		GetByTokenFunc: func(ctx context.Context, t string) (*models.PasswordResetToken, error) {
			return &models.PasswordResetToken{
				ID:        uuid.New(),
				UserID:    userID,
				Token:     token,
				ExpiresAt: time.Now().Add(1 * time.Hour),
				Used:      false,
			}, nil
		},
		MarkUsedFunc: func(ctx context.Context, id uuid.UUID, t time.Time) error {
			return nil
		},
	}
	userRepo := &MockUserRepository{
		UpdatePasswordFunc: func(ctx context.Context, id uuid.UUID, hash string) error {
			assert.Equal(t, userID, id)
			return nil
		},
	}

	svc := auth.NewAuthService(userRepo, prRepo, nil)
	err := svc.ResetPassword(ctx, token, "newsecurepassword123")
	assert.NoError(t, err)

	// Token expired case
	prRepo.GetByTokenFunc = func(ctx context.Context, t string) (*models.PasswordResetToken, error) {
		return &models.PasswordResetToken{
			UserID:    userID,
			ExpiresAt: time.Now().Add(-1 * time.Hour),
			Used:      false,
		}, nil
	}
	err = svc.ResetPassword(ctx, token, "newsecurepassword123")
	assert.ErrorContains(t, err, "token expired")
}

func TestAuthService_Refresh(t *testing.T) {
	t.Setenv("JWT_SECRET", "mysecret")
	t.Setenv("JWT_REFRESH_SECRET", "myrefreshsecret")

	ctx := context.Background()
	userID := uuid.New()
	email := "test@linkpulse.com"

	// Create valid refresh token
	refreshToken, err := auth.GenerateRefreshToken(userID, email, "testuser")
	require.NoError(t, err)

	userRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.User, error) {
			assert.Equal(t, userID, id)
			return &models.User{ID: userID, Email: email}, nil
		},
	}

	svc := auth.NewAuthService(userRepo, nil, nil)
	resp, err := svc.Refresh(ctx, refreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
}

func TestAuthService_LoginOrRegisterWithGoogle(t *testing.T) {
	t.Setenv("JWT_SECRET", "mysecret")
	t.Setenv("JWT_REFRESH_SECRET", "myrefreshsecret")

	ctx := context.Background()
	googleID := "g-12345"
	email := "google@linkpulse.com"
	name := "Google User"

	// Scenario 1: User already linked to Google ID
	userRepo := &MockUserRepository{
		GetByGoogleIDFunc: func(ctx context.Context, gID string) (*models.User, error) {
			return &models.User{ID: uuid.New(), Email: email, GoogleID: &googleID}, nil
		},
	}
	svc := auth.NewAuthService(userRepo, nil, nil)
	resp, err := svc.LoginOrRegisterWithGoogle(ctx, googleID, email, name)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)

	// Scenario 2: Register new user (not found by Google ID or Email)
	userID := uuid.New()
	userRepo = &MockUserRepository{
		GetByGoogleIDFunc: func(ctx context.Context, gID string) (*models.User, error) {
			return nil, nil
		},
		GetByEmailFunc: func(ctx context.Context, em string) (*models.User, error) {
			return nil, nil
		},
		GetByUsernameFunc: func(ctx context.Context, username string) (*models.User, error) {
			return nil, nil // Username is unique
		},
		CreateFunc: func(ctx context.Context, u *models.User) error {
			u.ID = userID
			return nil
		},
	}
	profileRepo := &MockProfileRepository{
		CreateProfileFunc: func(profile *models.Profile) error {
			assert.NotEmpty(t, profile.DisplayName)
			return nil
		},
	}
	svc = auth.NewAuthService(userRepo, nil, profileRepo)
	resp, err = svc.LoginOrRegisterWithGoogle(ctx, googleID, email, name)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
}

func TestAuthService_GetUserByID(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	userRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*models.User, error) {
			if id == userID {
				return &models.User{ID: userID, Email: "test@linkpulse.com"}, nil
			}
			return nil, nil
		},
	}
	svc := auth.NewAuthService(userRepo, nil, nil)
	user, err := svc.GetUserByID(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)

	_, err = svc.GetUserByID(ctx, uuid.New())
	assert.ErrorIs(t, err, auth.ErrUserNotFound)
}

// Helper hashing function for tests
func bcryptHash(password string) (string, error) {
	return utils.HashPassword(password)
}