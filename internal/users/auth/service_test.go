package auth_test

import (
	"context"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/authjwt"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/stretchr/testify/assert"
)

func setupTestService(t *testing.T) auth.AuthService {
	passwordsRepo := passwords.NewPasswordsRepo(storage)
	passwordsService := passwords.NewPasswordsService(passwordsRepo)
	authRepo := auth.NewAuthRepo(storage)

	t.Cleanup(func() {
		err := storage.DB.Exec("DELETE FROM passwords").Error
		assert.NoError(t, err)

		err = storage.Cache.FlushDB(context.Background()).Err()
		assert.NoError(t, err)
	})

	return auth.NewAuthService(authRepo, passwordsService)
}

func TestAuthService_Register(t *testing.T) {
	service := setupTestService(t)

	err := service.Register(1, "Password$123")
	assert.NoError(t, err)
}

func TestAuthService_Login(t *testing.T) {
	service := setupTestService(t)

	err := service.Register(1, "Password$123")
	assert.NoError(t, err)

	tokenPair, err := service.Login(1, "Password$123")
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.Access)
	assert.NotEmpty(t, tokenPair.Refresh)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	service := setupTestService(t)

	err := service.Register(1, "Password$123")
	assert.NoError(t, err)

	tokenPair, err := service.Login(1, "WrongPassword")
	assert.Error(t, err)
	assert.Nil(t, tokenPair)
}

func TestAuthService_Refresh(t *testing.T) {
	service := setupTestService(t)

	err := service.Register(1, "Password$123")
	assert.NoError(t, err)

	tokenPair, err := service.Login(1, "Password$123")
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)

	newTokenPair, err := service.Refresh(tokenPair.Refresh)
	assert.NoError(t, err)
	assert.NotNil(t, newTokenPair)
	assert.NotEmpty(t, newTokenPair.Access)
	assert.NotEmpty(t, newTokenPair.Refresh)
}

func TestAuthService_Refresh_InvalidToken(t *testing.T) {
	service := setupTestService(t)

	newTokenPair, err := service.Refresh("invalid token")
	assert.Error(t, err)
	assert.Nil(t, newTokenPair)
}

func TestAuthService_Logout(t *testing.T) {
	service := setupTestService(t)

	err := service.Register(1, "Password$123")
	assert.NoError(t, err)

	tokenPair, err := service.Login(1, "Password$123")
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)

	service.Logout(1)

	claims, err := service.IsTokenCached(&authjwt.JWTClaims{UserID: 1, Type: "access"})
	assert.NoError(t, err)
	assert.False(t, claims)
}
