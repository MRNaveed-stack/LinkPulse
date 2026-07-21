package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/MRNaveed-stack/LinkPulse/internal/auth"
	"github.com/MRNaveed-stack/LinkPulse/internal/handler"
	"github.com/MRNaveed-stack/LinkPulse/internal/middleware"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/repository"
	"github.com/MRNaveed-stack/LinkPulse/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkPulse_FullFlow_Integration(t *testing.T) {
	// 1. Setup Live DB Connection
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("Skipping integration test: TEST_DATABASE_URL not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	// Truncate tables to ensure tests run in clean isolation
	truncateTables(t, pool)

	// 2. Initialize Dependency Graph
	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewPasswordResetTokenRepository(pool)
	profileRepo := repository.NewProfileRepository(pool)
	linkRepo := repository.NewLinkRepository(pool)
	clickRepo := repository.NewClickRepository(pool)

	authSvc := auth.NewAuthService(userRepo, tokenRepo, profileRepo)
	linkSvc := service.NewLinkService(linkRepo, clickRepo)
	profileSvc := service.NewProfileService(userRepo, profileRepo, linkRepo)
	analyticsSvc := service.NewAnalyticsService(linkRepo)

	authHand := handler.NewAuthHandler(authSvc)
	linkHand := handler.NewLinkHandler(linkSvc, userRepo)
	profileHand := handler.NewProfileHandler(profileSvc)
	analyticsHand := handler.NewAnalyticsHandler(analyticsSvc)

	// Setup router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.Recovery())

	// Public routes
	r.POST("/api/v1/auth/register", authHand.Register)
	r.POST("/api/v1/auth/login", authHand.Login)
	r.GET("/r/:slug", linkHand.Redirect)
	r.GET("/api/v1/profiles/:username", profileHand.GetPublicProfile)

	// Authenticated routes
	authGroup := r.Group("/api/v1")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("/links", linkHand.CreateLink)
		authGroup.GET("/links", linkHand.GetUserLinks)
		authGroup.PUT("/profiles", profileHand.UpdateProfile)
		authGroup.GET("/analytics/overview", analyticsHand.GetOverview)
	}

	// Setup environmental variables for token validation
	t.Setenv("JWT_SECRET", "super-secret-integration-test-key")
	t.Setenv("JWT_REFRESH_SECRET", "super-refresh-secret-integration-test-key")

	// ==========================================
	// Step 1: User Registration
	// ==========================================
	w := httptest.NewRecorder()
	regReq := models.RegisterRequest{
		Username: "integrationuser",
		Email:    "integration@test.com",
		Password: "password12345",
	}
	body, _ := json.Marshal(regReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var regResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &regResp)
	require.NoError(t, err)

	accessToken := regResp["access_token"].(string)
	assert.NotEmpty(t, accessToken)

	// ==========================================
	// Step 2: User Login
	// ==========================================
	w = httptest.NewRecorder()
	loginReq := models.LoginRequest{
		Email:    "integration@test.com",
		Password: "password12345",
	}
	body, _ = json.Marshal(loginReq)
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// ==========================================
	// Step 3: Create Short Link
	// ==========================================
	w = httptest.NewRecorder()
	linkReq := models.CreateLinkRequest{
		Title:          "Golang",
		Slug:           "go",
		DestinationURL: "https://golang.org",
	}
	body, _ = json.Marshal(linkReq)
	req = httptest.NewRequest("POST", "/api/v1/links", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// ==========================================
	// Step 4: Handle Redirect & Track Click Event
	// ==========================================
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/r/go", nil)
	req.Header.Set("User-Agent", "Go-Integration-Test-Agent")
	req.Header.Set("Referer", "https://twitter.com/post")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "https://golang.org", w.Header().Get("Location"))

	// Give async DB writes a tiny breathing window
	time.Sleep(100 * time.Millisecond)

	// ==========================================
	// Step 5: Check Analytics overview
	// ==========================================
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/analytics/overview", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var overview models.AnalyticsOverview
	err = json.NewDecoder(w.Body).Decode(&overview)
	require.NoError(t, err)

	assert.Equal(t, int64(1), overview.TotalLinks)
	assert.Equal(t, int64(1), overview.TotalClicks)
	assert.Equal(t, int64(1), overview.ClicksToday)

	// ==========================================
	// Step 6: Create & View Profile
	// ==========================================
	w = httptest.NewRecorder()
	profileReq := models.UpdateProfileRequest{
		DisplayName: "Integration Master",
		Bio:         "Testing full flows.",
	}
	body, _ = json.Marshal(profileReq)
	req = httptest.NewRequest("PUT", "/api/v1/profiles", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Get public profile profile
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/api/v1/profiles/integrationuser", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var profileResp models.PublicProfileResponse
	err = json.NewDecoder(w.Body).Decode(&profileResp)
	require.NoError(t, err)
	assert.Equal(t, "Integration Master", profileResp.DisplayName)
	assert.Len(t, profileResp.Links, 1) // Active short link "go" should be visible
}

func truncateTables(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE click_events, password_reset_tokens, profiles, links, users CASCADE")
	require.NoError(t, err)
}
