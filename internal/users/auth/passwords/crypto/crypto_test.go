package crypto_test

import (
	"testing"

	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestCryptoService() crypto.CryptoService {
	return crypto.NewCryptoService()
}

func TestGenerateHash(t *testing.T) {
	service := setupTestCryptoService()
	password := []byte("mysecretpassword")

	hashedPassword, err := service.GenerateHash(password)
	require.NoError(t, err)
	require.NotNil(t, hashedPassword)

	assert.Equal(t, "argon2id", hashedPassword.Algorithm)
	assert.NotEmpty(t, hashedPassword.Hash)
	assert.NotEmpty(t, hashedPassword.Salt)
}

func TestCompare(t *testing.T) {
	service := setupTestCryptoService()
	password := []byte("mysecretpassword")

	hashedPassword, err := service.GenerateHash(password)
	require.NoError(t, err)
	require.NotNil(t, hashedPassword)

	match, err := service.Compare(password, hashedPassword)
	require.NoError(t, err)
	assert.True(t, match)

	wrongPassword := []byte("wrongpassword")
	match, err = service.Compare(wrongPassword, hashedPassword)
	require.NoError(t, err)
	assert.False(t, match)
}
