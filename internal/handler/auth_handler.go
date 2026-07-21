package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/MRNaveed-stack/LinkPulse/internal/auth"
	"github.com/MRNaveed-stack/LinkPulse/internal/logger"
	"github.com/MRNaveed-stack/LinkPulse/internal/models"
	"github.com/MRNaveed-stack/LinkPulse/internal/validator"
)

type AuthHandler struct {
	authService auth.AuthService
}

func NewAuthHandler(authService auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed JSON"})
		return
	}

	if err := validator.Validate(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}
	tokenResponse, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		log.Printf("Registration failed for user %s: %v", req.Username, err)
		if err == auth.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "User registered successfully",
		"access_token":  tokenResponse.AccessToken,
		"refresh_token": tokenResponse.RefreshToken,
		"token":         tokenResponse.Token,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	tokenResponse, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		log.Printf("Login failed for email %s: %v", req.Email, err)
		if err == auth.ErrInvalidEmail || err == auth.ErrInvalidPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

// ForgotPassword godoc
// @Summary Request password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.ForgotPasswordRequest true "Email address"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	token, err := h.authService.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		log.Printf("ForgotPassword failed for email %s: %v", req.Email, err)
		if err == auth.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process request"})
		return
	}

	// For MVP: Return the token in response
	// In production: Send via email
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset link sent to your email",
		"token":   token, // Only for development/testing
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":  validator.FormatValidationErrors(err),
			"details": err.Error(),
		})
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		log.Printf("ResetPassword failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	// Get user ID from context (set by middleware)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if err == auth.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user info"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func getGoogleOAuthConfig() *oauth2.Config {
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	// Enforce valid redirect URI as per requirements
	if redirectURL == "" || strings.Contains(redirectURL, "0.0.0.0") || strings.Contains(redirectURL, "/auth/google/login") {
		redirectURL = "http://localhost:8080/auth/google/callback"
	}
	// Clean trailing slash if present
	redirectURL = strings.TrimSuffix(redirectURL, "/")

	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GoogleLogin redirects user to Google OAuth consent page
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	config := getGoogleOAuthConfig()
	if config.ClientID == "" || config.ClientSecret == "" {
		log.Println("Google OAuth config is missing: check GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth is not configured on the server"})
		return
	}

	// Generate a secure random state
	stateBytes := make([]byte, 16)
	state := "state-token"
	if _, err := rand.Read(stateBytes); err == nil {
		state = hex.EncodeToString(stateBytes)
	}

	// Set state cookie (5-minute expiry)
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	url := config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	log.Printf("Google OAuth: Redirecting user to URL: %s", url)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the Google OAuth redirect redirecting back with authorization code
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	cookieState, err := c.Cookie("oauth_state")
	if err != nil || state == "" || state != cookieState {
		log.Printf("Google Callback: state validation failed. query state: %s, cookie state: %s, err: %v", state, cookieState, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
		return
	}
	// Clear the state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	code := c.Query("code")
	if code == "" {
		log.Println("Google Callback: missing auth code parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code parameter"})
		return
	}

	config := getGoogleOAuthConfig()
	token, err := config.Exchange(c.Request.Context(), code)
	if err != nil {
		log.Printf("Google Callback: token exchange failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange authorization code: " + err.Error()})
		return
	}

	userInfo, err := fetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		log.Printf("Google Callback: failed to fetch user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info from Google: " + err.Error()})
		return
	}

	tokenResponse, err := h.authService.LoginOrRegisterWithGoogle(c.Request.Context(), userInfo.ID, userInfo.Email, userInfo.Name)
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login?error="+err.Error())
		return
	}

	redirectTarget := fmt.Sprintf(
		"%s/login?access_token=%s&refresh_token=%s",
		frontendURL, tokenResponse.AccessToken, tokenResponse.RefreshToken,
	)
	c.Redirect(http.StatusTemporaryRedirect, redirectTarget)
}

type GoogleTokenRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

// GoogleTokenExchange validates a Google Access Token directly and issues a JWT token
func (h *AuthHandler) GoogleTokenExchange(c *gin.Context) {
	var req GoogleTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: access_token is required"})
		return
	}

	userInfo, err := fetchGoogleUserInfo(req.AccessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to verify token with Google: " + err.Error()})
		return
	}

	tokenResponse, err := h.authService.LoginOrRegisterWithGoogle(c.Request.Context(), userInfo.ID, userInfo.Email, userInfo.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  tokenResponse.AccessToken,
		"refresh_token": tokenResponse.RefreshToken,
		"token":         tokenResponse.Token,
	})
}

type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func fetchGoogleUserInfo(accessToken string) (*googleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google returned non-OK status: %d", resp.StatusCode)
	}

	var info googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	if info.Email == "" {
		return nil, fmt.Errorf("email not provided by Google")
	}

	return &info, nil
}

// Refresh godoc
// @Summary Refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RefreshRequest true "Refresh token"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Warn(
			"validation failed",
			"endpoint", c.FullPath(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request: refresh_token is required",
			"details": err.Error(),
		})
		return
	}

	tokenResponse, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}
