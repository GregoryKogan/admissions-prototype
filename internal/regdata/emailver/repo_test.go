package emailver_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/regdata/emailver"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	storage datastore.MockStorage
)

func TestMain(m *testing.M) {
	s, cleanup := datastore.InitMockStorage()
	storage = s

	viper.Set("auth.email_verification.token_lifetime", "15m")

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func setupTestRepo(t *testing.T) emailver.EmailVerificationRepo {
	t.Cleanup(func() {
		err := storage.Flush()
		assert.NoError(t, err)
	})

	return emailver.NewEmailVerificationRepo(storage)
}

func TestCreateVerificationToken(t *testing.T) {
	repo := setupTestRepo(t)

	token, err := repo.CreateVerificationToken(1)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token in Redis
	val, err := storage.Cache().Get(context.Background(), fmt.Sprintf("email-token:%s", token)).Result()
	require.NoError(t, err)
	assert.Equal(t, "1", val)

	// Verify TTL was set
	ttl, err := storage.Cache().TTL(context.Background(), fmt.Sprintf("email-token:%s", token)).Result()
	require.NoError(t, err)
	assert.True(t, ttl > 0)
}

func TestGetRegistrationIDByToken(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	// Test with invalid token
	_, err := repo.GetRegistrationIDByToken("invalid-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired token")

	// Test with valid token
	token, err := repo.CreateVerificationToken(1)
	require.NoError(t, err)

	registrationID, err := repo.GetRegistrationIDByToken(token)
	require.NoError(t, err)
	assert.Equal(t, uint(1), registrationID)

	// Test expired token
	err = storage.Cache().Del(ctx, fmt.Sprintf("email-token:%s", token)).Err()
	require.NoError(t, err)

	_, err = repo.GetRegistrationIDByToken(token)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired token")
}

func TestDeleteToken(t *testing.T) {
	repo := setupTestRepo(t)
	ctx := context.Background()

	// Create and verify token exists
	token, err := repo.CreateVerificationToken(1)
	require.NoError(t, err)

	// Delete token
	err = repo.DeleteToken(token)
	require.NoError(t, err)

	// Verify deletion
	exists, err := storage.Cache().Exists(ctx, fmt.Sprintf("email-token:%s", token)).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)

	// Test deleting non-existent token
	err = repo.DeleteToken("non-existent-token")
	assert.NoError(t, err) // Should not error when deleting non-existent key
}
