package emailver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/regdata/emailver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestService(t *testing.T) emailver.EmailVerificationService {
	t.Cleanup(func() {
		err := storage.Flush()
		assert.NoError(t, err)
	})

	repo := emailver.NewEmailVerificationRepo(storage)
	return emailver.NewEmailVerificationService(repo)
}

func getTokenFromRedis(t *testing.T, registrationID uint) string {
	pattern := "email-token:*"
	keys, err := storage.Cache().Keys(context.Background(), pattern).Result()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	// Verify the token maps to correct registration ID
	val, err := storage.Cache().Get(context.Background(), keys[0]).Result()
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprint(registrationID), val)

	return keys[0][len("email-token:"):]
}

func TestSendVerificationEmail(t *testing.T) {
	service := setupTestService(t)

	err := service.SendVerificationEmail("test@example.com", 1)
	assert.NoError(t, err)

	// Verify token was stored in Redis
	token := getTokenFromRedis(t, 1)
	assert.NotEmpty(t, token)
}

func TestVerifyEmail(t *testing.T) {
	service := setupTestService(t)

	// Create verification token
	err := service.SendVerificationEmail("test@example.com", 1)
	require.NoError(t, err)
	token := getTokenFromRedis(t, 1)

	// Test verification with valid token
	registrationID, err := service.VerifyEmail(token)
	require.NoError(t, err)
	assert.Equal(t, uint(1), registrationID)

	// Verify token was deleted from Redis
	exists, err := storage.Cache().Exists(context.Background(), fmt.Sprintf("email-token:%s", token)).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), exists)

	// Test verification with invalid token
	_, err = service.VerifyEmail("invalid-token")
	assert.Error(t, err)
}
