package handler_test

import (
	"context"

	"github.com/MRNaveed-stack/LinkPulse/internal/auth"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/google/uuid"
)

type MockAuthService struct {
	RegisterFunc                  func(ctx context.Context, req models.RegisterRequest) (*models.LoginResponse, error)
	LoginFunc                     func(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
	ForgotPasswordFunc            func(ctx context.Context, email string) (string, error)
	ResetPasswordFunc             func(ctx context.Context, token string, newPassword string) error
	ValidateTokenFunc             func(token string) (*auth.TokenClaims, error)
	LoginOrRegisterWithGoogleFunc func(ctx context.Context, googleID, email, name string) (*models.LoginResponse, error)
	RefreshFunc                   func(ctx context.Context, refreshToken string) (*models.LoginResponse, error)
	GetUserByIDFunc               func(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

func (m *MockAuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.LoginResponse, error) {
	return m.RegisterFunc(ctx, req)
}
func (m *MockAuthService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	return m.LoginFunc(ctx, req)
}
func (m *MockAuthService) ForgotPassword(ctx context.Context, email string) (string, error) {
	return m.ForgotPasswordFunc(ctx, email)
}
func (m *MockAuthService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	return m.ResetPasswordFunc(ctx, token, newPassword)
}
func (m *MockAuthService) ValidateToken(token string) (*auth.TokenClaims, error) {
	return m.ValidateTokenFunc(token)
}
func (m *MockAuthService) LoginOrRegisterWithGoogle(ctx context.Context, googleID, email, name string) (*models.LoginResponse, error) {
	return m.LoginOrRegisterWithGoogleFunc(ctx, googleID, email, name)
}
func (m *MockAuthService) Refresh(ctx context.Context, refreshToken string) (*models.LoginResponse, error) {
	return m.RefreshFunc(ctx, refreshToken)
}
func (m *MockAuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return m.GetUserByIDFunc(ctx, userID)
}

type MockLinkService struct {
	CreateLinkFunc       func(userID uuid.UUID, req models.CreateLinkRequest) error
	GetUserLinksFunc     func(userID uuid.UUID) ([]*models.Link, error)
	HandleRedirectFunc   func(userID uuid.UUID, slug string, ip string, userAgent string, referrer string) (string, error)
	UpdateLinkFunc       func(userID uuid.UUID, linkID uuid.UUID, req models.UpdateLinkRequest) (*models.Link, error)
	DeleteLinkFunc       func(userID uuid.UUID, linkID uuid.UUID) error
	UpdateLinkStatusFunc func(userID uuid.UUID, linkID uuid.UUID, isActive bool) error
}

func (m *MockLinkService) CreateLink(userID uuid.UUID, req models.CreateLinkRequest) error {
	return m.CreateLinkFunc(userID, req)
}
func (m *MockLinkService) GetUserLinks(userID uuid.UUID) ([]*models.Link, error) {
	return m.GetUserLinksFunc(userID)
}
func (m *MockLinkService) HandleRedirect(userID uuid.UUID, slug string, ip string, userAgent string, referrer string) (string, error) {
	return m.HandleRedirectFunc(userID, slug, ip, userAgent, referrer)
}
func (m *MockLinkService) UpdateLink(userID uuid.UUID, linkID uuid.UUID, req models.UpdateLinkRequest) (*models.Link, error) {
	return m.UpdateLinkFunc(userID, linkID, req)
}
func (m *MockLinkService) DeleteLink(userID uuid.UUID, linkID uuid.UUID) error {
	return m.DeleteLinkFunc(userID, linkID)
}
func (m *MockLinkService) UpdateLinkStatus(userID uuid.UUID, linkID uuid.UUID, isActive bool) error {
	return m.UpdateLinkStatusFunc(userID, linkID, isActive)
}

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
	return m.UpdateUsernameFunc(ctx, id, username)
}
func (m *MockUserRepository) UpdateEmail(ctx context.Context, id uuid.UUID, email string) error {
	return m.UpdateEmailFunc(ctx, id, email)
}
func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.DeleteFunc(ctx, id)
}

type MockProfileService struct {
	GetPublicProfileFunc func(username string) (*models.PublicProfileResponse, error)
	UpdateProfileFunc    func(userID uuid.UUID, req *models.UpdateProfileRequest) (*models.Profile, error)
	ChangePasswordFunc   func(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error
	DeleteAccountFunc    func(ctx context.Context, userID uuid.UUID) error
}

func (m *MockProfileService) GetPublicProfile(username string) (*models.PublicProfileResponse, error) {
	return m.GetPublicProfileFunc(username)
}
func (m *MockProfileService) UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) (*models.Profile, error) {
	return m.UpdateProfileFunc(userID, req)
}
func (m *MockProfileService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	if m.ChangePasswordFunc != nil {
		return m.ChangePasswordFunc(ctx, userID, currentPassword, newPassword)
	}
	return nil
}
func (m *MockProfileService) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	if m.DeleteAccountFunc != nil {
		return m.DeleteAccountFunc(ctx, userID)
	}
	return nil
}

type MockAnalyticsService struct {
	GetOverviewFunc          func(userID uuid.UUID) (*models.AnalyticsOverview, error)
	GetLinkAnalyticsFunc     func(userID uuid.UUID) ([]models.LinkAnalyticsDTO, error)
	GetDailyAnalyticsFunc    func(userID uuid.UUID, days int) ([]models.DailyAnalyticsDTO, error)
	GetRecentActivityFunc    func(userID uuid.UUID, limit int) ([]models.RecentActivityDTO, error)
	GetReferrerAnalyticsFunc func(userID uuid.UUID) ([]models.ReferrerAnalyticsDTO, error)
}

func (m *MockAnalyticsService) GetOverview(userID uuid.UUID) (*models.AnalyticsOverview, error) {
	return m.GetOverviewFunc(userID)
}
func (m *MockAnalyticsService) GetLinkAnalytics(userID uuid.UUID) ([]models.LinkAnalyticsDTO, error) {
	return m.GetLinkAnalyticsFunc(userID)
}
func (m *MockAnalyticsService) GetDailyAnalytics(userID uuid.UUID, days int) ([]models.DailyAnalyticsDTO, error) {
	return m.GetDailyAnalyticsFunc(userID, days)
}
func (m *MockAnalyticsService) GetRecentActivity(userID uuid.UUID, limit int) ([]models.RecentActivityDTO, error) {
	return m.GetRecentActivityFunc(userID, limit)
}
func (m *MockAnalyticsService) GetReferrerAnalytics(userID uuid.UUID) ([]models.ReferrerAnalyticsDTO, error) {
	return m.GetReferrerAnalyticsFunc(userID)
}
