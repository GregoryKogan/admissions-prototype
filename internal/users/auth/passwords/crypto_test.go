package passwords_test

import (
	"testing"

	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateHash(t *testing.T) {
	params := passwords.DefaultArgon2idParams()
	password := []byte("mysecretpassword")

	hashedPassword, err := params.GenerateHash(password)
	require.NoError(t, err)
	require.NotNil(t, hashedPassword)

	assert.Equal(t, "argon2id", hashedPassword.Algorithm)
	assert.NotEmpty(t, hashedPassword.Hash)
	assert.NotEmpty(t, hashedPassword.Salt)
}

func TestCompare(t *testing.T) {
	params := passwords.DefaultArgon2idParams()
	password := []byte("mysecretpassword")

	hashedPassword, err := params.GenerateHash(password)
	require.NoError(t, err)
	require.NotNil(t, hashedPassword)

	match, err := params.Compare(password, hashedPassword)
	require.NoError(t, err)
	assert.True(t, match)

	wrongPassword := []byte("wrongpassword")
	match, err = params.Compare(wrongPassword, hashedPassword)
	require.NoError(t, err)
	assert.False(t, match)
}
