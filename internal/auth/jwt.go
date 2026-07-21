package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Type     string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT access token (backward compatible wrapper)
func GenerateToken(userID uuid.UUID, email string, username string) (string, error) {
	return GenerateAccessToken(userID, email, username)
}

// GenerateTokenPair generates both an access token and a refresh token
func GenerateTokenPair(userID uuid.UUID, email string, username string) (string, string, error) {
	accessToken, err := GenerateAccessToken(userID, email, username)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := GenerateRefreshToken(userID, email, username)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// GenerateAccessToken creates a new JWT access token
func GenerateAccessToken(userID uuid.UUID, email string, username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set in environment")
	}

	expirationHours := os.Getenv("JWT_EXPIRATION_HOURS")
	if expirationHours == "" {
		expirationHours = "24" // Default 24 hours
	}

	expiration, err := time.ParseDuration(expirationHours + "h")
	if err != nil {
		return "", fmt.Errorf("invalid JWT_EXPIRATION_HOURS: %w", err)
	}

	claims := Claims{
		UserID:   userID.String(),
		Email:    email,
		Username: username,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "LinkPulse",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a new JWT refresh token
func GenerateRefreshToken(userID uuid.UUID, email string, username string) (string, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = os.Getenv("JWT_SECRET")
	}
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET not set in environment")
	}

	expirationDays := os.Getenv("JWT_REFRESH_EXPIRATION_DAYS")
	days := 30
	if expirationDays != "" {
		if val, err := strconv.Atoi(expirationDays); err == nil {
			days = val
		}
	}
	expiration := time.Duration(days) * 24 * time.Hour

	claims := Claims{
		UserID:   userID.String(),
		Email:    email,
		Username: username,
		Type:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "LinkPulse",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates and parses a JWT token (access token)
func ValidateToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set in environment")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// ValidateRefreshToken parses and validates a refresh token
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = os.Getenv("JWT_SECRET")
	}
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set in environment")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.Type != "refresh" {
			return nil, fmt.Errorf("invalid token type: expected refresh")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
