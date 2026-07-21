package auth_test

import (
	"context"
	"time"

	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/google/uuid"
)

// MockUserRepository mocks repository.UserRepository
type MockUserRepository struct {
	CreateFunc            func(ctx context.Context, user *models.User) error
	GetByEmailFunc        func(ctx context.Context, email string) (*models.User, error)
	GetByIDFunc           func(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByUsernameFunc     func(ctx context.Context, username string) (*models.User, error)
	GetByGoogleIDFunc     func(ctx context.Context, googleID string) (*models.User, error)
	LinkGoogleAccountFunc func(ctx context.Context, id uuid.UUID, googleID string) error
	UpdatePasswordFunc    func(ctx context.Context, id uuid.UUID, hash string) error
	UpdatePlanFunc        func(ctx context.Context, id uuid.UUID, plan string) error
	UpdateUsernameFunc    func(ctx context.Context, id uuid.UUID, username string) error
	UpdateEmailFunc       func(ctx context.Context, id uuid.UUID, email string) error
	DeleteFunc            func(ctx context.Context, id uuid.UUID) error
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	return m.CreateFunc(ctx, user)
}
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.GetByEmailFunc(ctx, email)
}
func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return m.GetByUsernameFunc(ctx, username)
}
func (m *MockUserRepository) GetByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	return m.GetByGoogleIDFunc(ctx, googleID)
}
func (m *MockUserRepository) LinkGoogleAccount(ctx context.Context, id uuid.UUID, googleID string) error {
	return m.LinkGoogleAccountFunc(ctx, id, googleID)
}
func (m *MockUserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error {
	return m.UpdatePasswordFunc(ctx, id, hash)
}
func (m *MockUserRepository) UpdatePlan(ctx context.Context, id uuid.UUID, plan string) error {
	return m.UpdatePlanFunc(ctx, id, plan)
}
func (m *MockUserRepository) UpdateUsername(ctx context.Context, id uuid.UUID, username string) error {
	if m.UpdateUsernameFunc != nil {
		return m.UpdateUsernameFunc(ctx, id, username)
	}
	return nil
}
func (m *MockUserRepository) UpdateEmail(ctx context.Context, id uuid.UUID, email string) error {
	if m.UpdateEmailFunc != nil {
		return m.UpdateEmailFunc(ctx, id, email)
	}
	return nil
}
func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// MockPasswordResetTokenRepository mocks repository.PasswordResetTokenRepository
type MockPasswordResetTokenRepository struct {
	CreateFunc     func(ctx context.Context, token *models.PasswordResetToken) error
	GetByTokenFunc func(ctx context.Context, token string) (*models.PasswordResetToken, error)
	MarkUsedFunc   func(ctx context.Context, tokenID uuid.UUID, usedAt time.Time) error
}

func (m *MockPasswordResetTokenRepository) Create(ctx context.Context, token *models.PasswordResetToken) error {
	return m.CreateFunc(ctx, token)
}
func (m *MockPasswordResetTokenRepository) GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error) {
	return m.GetByTokenFunc(ctx, token)
}
func (m *MockPasswordResetTokenRepository) MarkUsed(ctx context.Context, tokenID uuid.UUID, usedAt time.Time) error {
	return m.MarkUsedFunc(ctx, tokenID, usedAt)
}

// MockProfileRepository mocks repository.ProfileRepository
type MockProfileRepository struct {
	GetByUserIDFunc   func(userID uuid.UUID) (*models.Profile, error)
	CreateProfileFunc func(profile *models.Profile) error
	UpdateProfileFunc func(profile *models.Profile) error
}

func (m *MockProfileRepository) GetByUserID(userID uuid.UUID) (*models.Profile, error) {
	return m.GetByUserIDFunc(userID)
}
func (m *MockProfileRepository) CreateProfile(profile *models.Profile) error {
	return m.CreateProfileFunc(profile)
}
func (m *MockProfileRepository) UpdateProfile(profile *models.Profile) error {
	return m.UpdateProfileFunc(profile)
}
