package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/auth"
	"github.com/MRNaveed-stack/LinkPulse/internal/handler"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := &MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	// Case 1: Valid Register request
	mockSvc.RegisterFunc = func(ctx context.Context, req models.RegisterRequest) (*models.LoginResponse, error) {
		return &models.LoginResponse{AccessToken: "token123"}, nil
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.POST("/register", h.Register)

	reqBody, _ := json.Marshal(models.RegisterRequest{
		Username: "newuser",
		Email:    "test@linkpulse.com",
		Password: "password123",
	})
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Case 2: Validation Failure
	w = httptest.NewRecorder()
	c, r = gin.CreateTestContext(w)
	r.POST("/register", h.Register)

	reqBody, _ = json.Marshal(models.RegisterRequest{
		Username: "nu",
		Email:    "not-an-email",
		Password: "short",
	})
	req = httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	c.Request = req
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	mockSvc.LoginFunc = func(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
		if req.Email == "correct@email.com" && req.Password == "correct" {
			return &models.LoginResponse{AccessToken: "access"}, nil
		}
		return nil, auth.ErrInvalidPassword
	}

	// Success case
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	r.POST("/login", h.Login)

	reqBody, _ := json.Marshal(models.LoginRequest{Email: "correct@email.com", Password: "correct"})
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Failure case
	w2 := httptest.NewRecorder()
	reqBody2, _ := json.Marshal(models.LoginRequest{Email: "correct@email.com", Password: "wrong"})
	req2 := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody2))
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code)
}

func TestAuthHandler_ForgotPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	mockSvc.ForgotPasswordFunc = func(ctx context.Context, email string) (string, error) {
		if email == "exist@linkpulse.com" {
			return "resettoken", nil
		}
		return "", auth.ErrUserNotFound
	}

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(models.ForgotPasswordRequest{Email: "exist@linkpulse.com"})
	req := httptest.NewRequest("POST", "/forgot", bytes.NewBuffer(reqBody))

	_, r := gin.CreateTestContext(w)
	r.POST("/forgot", h.ForgotPassword)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	mockSvc.ResetPasswordFunc = func(ctx context.Context, token, pass string) error {
		return nil
	}

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(models.ResetPasswordRequest{Token: "token", NewPassword: "newsecurepass"})
	req := httptest.NewRequest("POST", "/reset", bytes.NewBuffer(reqBody))

	_, r := gin.CreateTestContext(w)
	r.POST("/reset", h.ResetPassword)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_GetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	userID := uuid.New()
	mockSvc.GetUserByIDFunc = func(ctx context.Context, id uuid.UUID) (*models.User, error) {
		return &models.User{ID: id, Username: "me"}, nil
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/me", nil)

	_, r := gin.CreateTestContext(w)
	r.GET("/me", func(c *gin.Context) {
		c.Set("user_id", userID)
	}, h.GetMe)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_GoogleLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/login/google", nil)

	_, r := gin.CreateTestContext(w)
	r.GET("/login/google", h.GoogleLogin)
	r.ServeHTTP(w, req)

	// Since OAuth credentials are not set, it should log or redirect depending on setup.
	assert.Contains(t, []int{http.StatusTemporaryRedirect, http.StatusInternalServerError}, w.Code)
}

func TestAuthHandler_Refresh(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &MockAuthService{}
	h := handler.NewAuthHandler(mockSvc)

	mockSvc.RefreshFunc = func(ctx context.Context, rt string) (*models.LoginResponse, error) {
		return &models.LoginResponse{AccessToken: "new-access"}, nil
	}

	w := httptest.NewRecorder()
	reqBody, _ := json.Marshal(models.RefreshRequest{RefreshToken: "refresh-token"})
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(reqBody))

	_, r := gin.CreateTestContext(w)
	r.POST("/refresh", h.Refresh)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
