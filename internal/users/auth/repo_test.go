package auth_test

import (
	"os"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/authjwt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	storage datastore.MockStorage
)

func TestMain(m *testing.M) {
	s, cleanup := datastore.InitMockStorage()
	storage = s

	viper.Set("auth.access_lifetime", "15m")
	viper.Set("auth.refresh_lifetime", "720h")
	viper.Set("auth.auto_logout", "24h")
	viper.Set("secrets.jwt_key", "test_key")

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func setupTestRepo(t *testing.T) auth.AuthRepo {
	t.Cleanup(func() {
		err := storage.Flush()
		assert.NoError(t, err)
	})

	return auth.NewAuthRepo(storage)
}

func TestCacheTokenPair(t *testing.T) {
	repo := setupTestRepo(t)

	jwtService := authjwt.NewJWTService()

	access, err := jwtService.NewAccessToken(1)
	assert.NoError(t, err)

	refresh, err := jwtService.NewRefreshToken(1)
	assert.NoError(t, err)

	pair := &auth.TokenPair{
		Access:  access,
		Refresh: refresh,
	}

	err = repo.CacheTokenPair(pair)
	assert.NoError(t, err)
}

func TestCacheTokenPair_InvalidTokenPair(t *testing.T) {
	repo := setupTestRepo(t)

	jwtService := authjwt.NewJWTService()

	access, err := jwtService.NewAccessToken(1)
	assert.NoError(t, err)

	refresh, err := jwtService.NewRefreshToken(2)
	assert.NoError(t, err)

	pair := &auth.TokenPair{
		Access:  access,
		Refresh: refresh,
	}

	err = repo.CacheTokenPair(pair)
	assert.Error(t, err)
}

func TestCacheTokenPair_InvalidToken(t *testing.T) {
	repo := setupTestRepo(t)

	pair := &auth.TokenPair{
		Access:  "invalid token",
		Refresh: "invalid token",
	}

	err := repo.CacheTokenPair(pair)
	assert.Error(t, err)
}

func TestIsTokenCached(t *testing.T) {
	repo := setupTestRepo(t)

	jwtService := authjwt.NewJWTService()

	access, err := jwtService.NewAccessToken(1)
	assert.NoError(t, err)

	refresh, err := jwtService.NewRefreshToken(1)
	assert.NoError(t, err)

	pair := &auth.TokenPair{
		Access:  access,
		Refresh: refresh,
	}

	err = repo.CacheTokenPair(pair)
	assert.NoError(t, err)

	accessClaims, err := jwtService.ParseToken(access)
	assert.NoError(t, err)
	cached, err := repo.IsTokenCached(accessClaims)
	assert.NoError(t, err)
	assert.True(t, cached)

	refreshClaims, err := jwtService.ParseToken(refresh)
	assert.NoError(t, err)
	cached, err = repo.IsTokenCached(refreshClaims)
	assert.NoError(t, err)
	assert.True(t, cached)
}

func TestIsTokenCached_NotCached(t *testing.T) {
	repo := setupTestRepo(t)

	jwtService := authjwt.NewJWTService()

	access, err := jwtService.NewAccessToken(1)
	assert.NoError(t, err)

	accessClaims, err := jwtService.ParseToken(access)
	assert.NoError(t, err)
	cached, err := repo.IsTokenCached(accessClaims)
	assert.NoError(t, err)
	assert.False(t, cached)
}

func TestDeleteTokenPair(t *testing.T) {
	repo := setupTestRepo(t)

	jwtService := authjwt.NewJWTService()

	access, err := jwtService.NewAccessToken(1)
	assert.NoError(t, err)

	refresh, err := jwtService.NewRefreshToken(1)
	assert.NoError(t, err)

	pair := &auth.TokenPair{
		Access:  access,
		Refresh: refresh,
	}

	err = repo.CacheTokenPair(pair)
	assert.NoError(t, err)

	// Ensure tokens are cached
	accessClaims, err := jwtService.ParseToken(access)
	assert.NoError(t, err)
	cached, err := repo.IsTokenCached(accessClaims)
	assert.NoError(t, err)
	assert.True(t, cached)

	refreshClaims, err := jwtService.ParseToken(refresh)
	assert.NoError(t, err)
	cached, err = repo.IsTokenCached(refreshClaims)
	assert.NoError(t, err)
	assert.True(t, cached)

	// Delete token pair
	repo.DeleteTokenPair(1)

	// Ensure tokens are no longer cached
	cached, err = repo.IsTokenCached(accessClaims)
	assert.NoError(t, err)
	assert.False(t, cached)

	cached, err = repo.IsTokenCached(refreshClaims)
	assert.NoError(t, err)
	assert.False(t, cached)
}

func TestExtendTokenPairCacheExpiration(t *testing.T) {
	repo := setupTestRepo(t)

	jwtService := authjwt.NewJWTService()

	access, err := jwtService.NewAccessToken(1)
	assert.NoError(t, err)

	refresh, err := jwtService.NewRefreshToken(1)
	assert.NoError(t, err)

	pair := &auth.TokenPair{
		Access:  access,
		Refresh: refresh,
	}

	err = repo.CacheTokenPair(pair)
	assert.NoError(t, err)

	// Ensure tokens are cached
	accessClaims, err := jwtService.ParseToken(access)
	assert.NoError(t, err)
	cached, err := repo.IsTokenCached(accessClaims)
	assert.NoError(t, err)
	assert.True(t, cached)

	refreshClaims, err := jwtService.ParseToken(refresh)
	assert.NoError(t, err)
	cached, err = repo.IsTokenCached(refreshClaims)
	assert.NoError(t, err)
	assert.True(t, cached)

	// Extend token pair cache expiration
	repo.ExtendTokenPairCacheExpiration(1)

	// Ensure tokens are still cached after expiration extension
	cached, err = repo.IsTokenCached(accessClaims)
	assert.NoError(t, err)
	assert.True(t, cached)

	cached, err = repo.IsTokenCached(refreshClaims)
	assert.NoError(t, err)
	assert.True(t, cached)
}
