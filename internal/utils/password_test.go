package utils_test

import (
	"testing"

	"github.com/MRNaveed-stack/LinkPulse/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHelpers(t *testing.T) {
	password := "supersecure123"

	// 1. Test HashPassword
	hash, err := utils.HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// 2. Test CheckPassword - Success Case
	err = utils.CheckPassword(password, hash)
	assert.NoError(t, err)

	// 3. Test CheckPassword - Failure Case
	err = utils.CheckPassword("wrongpassword", hash)
	assert.Error(t, err)
	assert.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
}
