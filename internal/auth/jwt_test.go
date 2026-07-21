package auth_test

import (
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	// Set environment variables for tests
	t.Setenv("JWT_SECRET", "super-secret-key-for-testing-12345")
	t.Setenv("JWT_REFRESH_SECRET", "refresh-super-secret-key-for-testing-12345")
	t.Setenv("JWT_EXPIRATION_HOURS", "1")
	t.Setenv("JWT_REFRESH_EXPIRATION_DAYS", "7")

	userID := uuid.New()
	email := "test@linkpulse.com"

	// 1. Generate Token Pair
	accessToken, refreshToken, err := auth.GenerateTokenPair(userID, email, "testuser")
	require.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	// 2. Validate Access Token
	accessClaims, err := auth.ValidateToken(accessToken)
	require.NoError(t, err)
	assert.Equal(t, userID.String(), accessClaims.UserID)
	assert.Equal(t, email, accessClaims.Email)
	assert.Equal(t, "testuser", accessClaims.Username)
	assert.Equal(t, userID.String(), accessClaims.Subject)
	assert.Equal(t, "access", accessClaims.Type)

	// 3. Validate Refresh Token
	refreshClaims, err := auth.ValidateRefreshToken(refreshToken)
	require.NoError(t, err)
	assert.Equal(t, userID.String(), refreshClaims.UserID)
	assert.Equal(t, email, refreshClaims.Email)
	assert.Equal(t, "testuser", refreshClaims.Username)
	assert.Equal(t, userID.String(), refreshClaims.Subject)
	assert.Equal(t, "refresh", refreshClaims.Type)

	// 4. Reject Invalid Access Token
	_, err = auth.ValidateToken("invalid-token-string")
	assert.Error(t, err)

	// 5. Reject Invalid Refresh Token
	_, err = auth.ValidateRefreshToken("invalid-token-string")
	assert.Error(t, err)
}

func TestJWT_MissingSecrets(t *testing.T) {
	// Clear environments
	t.Setenv("JWT_SECRET", "")
	t.Setenv("JWT_REFRESH_SECRET", "")

	userID := uuid.New()
	email := "test@linkpulse.com"

	_, _, err := auth.GenerateTokenPair(userID, email, "testuser")
	assert.Error(t, err)
}
