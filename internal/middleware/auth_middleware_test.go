package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/auth"
	"github.com/MRNaveed-stack/LinkPulse/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	// Set mock JWT secret
	t.Setenv("JWT_SECRET", "test-secret-key-1234567890-abcdef")
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	email := "user@linkpulse.com"
	validToken, err := auth.GenerateAccessToken(userID, email, "testuser")
	require.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		verifyContext  bool
	}{
		{
			name:           "Valid Token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			verifyContext:  true,
		},
		{
			name:           "Missing Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			verifyContext:  false,
		},
		{
			name:           "Malformed Authorization Header",
			authHeader:     "Basic " + validToken,
			expectedStatus: http.StatusUnauthorized,
			verifyContext:  false,
		},
		{
			name:           "Invalid Signature Token",
			authHeader:     "Bearer " + validToken + "invalidstuff",
			expectedStatus: http.StatusUnauthorized,
			verifyContext:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.Use(middleware.AuthMiddleware())
			r.GET("/protected", func(ctx *gin.Context) {
				if tt.verifyContext {
					ctxID, existsID := ctx.Get("user_id")
					ctxEmail, existsEmail := ctx.Get("email")
					assert.True(t, existsID)
					assert.True(t, existsEmail)
					assert.Equal(t, userID, ctxID)
					assert.Equal(t, email, ctxEmail)
				}
				ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
			})

			req := httptest.NewRequest("GET", "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			c.Request = req
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
