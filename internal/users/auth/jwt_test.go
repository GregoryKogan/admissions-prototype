package auth_test

import (
	"testing"
	"time"

	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewAccessToken(t *testing.T) {
	viper.Set("jwt.access_lifetime", time.Hour)
	secrets.SetMockSecret("jwt_key", "testkey")
	defer secrets.ClearMockSecrets()

	tokenString, err := auth.NewAccessToken(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestNewRefreshToken(t *testing.T) {
	viper.Set("jwt.refresh_lifetime", time.Hour*24)
	secrets.SetMockSecret("jwt_key", "testkey")
	defer secrets.ClearMockSecrets()

	tokenString, err := auth.NewRefreshToken(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestParseToken(t *testing.T) {
	viper.Set("jwt.access_lifetime", time.Hour)
	secrets.SetMockSecret("jwt_key", "testkey")
	defer secrets.ClearMockSecrets()

	tokenString, err := auth.NewAccessToken(1)
	assert.NoError(t, err)

	claims, err := auth.ParseToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, "access", claims.Type)
}

func TestParseToken_Invalid(t *testing.T) {
	secrets.SetMockSecret("jwt_key", "testkey")
	defer secrets.ClearMockSecrets()

	_, err := auth.ParseToken("invalid.token")
	assert.Error(t, err)
}
