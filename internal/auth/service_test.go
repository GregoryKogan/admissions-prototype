package auth_test

import (
	"testing"

	"github.com/L2SH-Dev/admissions/internal/auth"
	"github.com/L2SH-Dev/admissions/internal/passwords"
	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/L2SH-Dev/admissions/internal/storage"
	"github.com/stretchr/testify/assert"
)

func setupTestService(t *testing.T) auth.AuthService {
	storage := storage.SetupMockStorage(t)
	passwordsRepo := passwords.NewPasswordsRepo(storage)
	passwordsService := passwords.NewPasswordsService(passwordsRepo)
	return auth.NewAuthService(passwordsService)
}

func TestAuthService_Register_And_Login(t *testing.T) {
	authService := setupTestService(t)

	secrets.SetMockSecret("jwt_key", "secret")
	defer secrets.ClearMockSecrets()

	// Register user
	err := authService.Register(1, "Password-123")
	assert.NoError(t, err)

	// Successful login
	tokenPair, err := authService.Login(1, "Password-123")
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	authService := setupTestService(t)

	// Register user
	err := authService.Register(1, "Password-123")
	assert.NoError(t, err)

	// Attempt login with wrong password
	tokenPair, err := authService.Login(1, "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, auth.ErrInvalidPassword, err)
	assert.Nil(t, tokenPair)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	authService := setupTestService(t)

	// Attempt login without registering
	tokenPair, err := authService.Login(1, "Password-123")
	assert.Error(t, err)
	assert.Nil(t, tokenPair)
}

func TestAuthService_ValidatePassword(t *testing.T) {
	authService := setupTestService(t)

	err := authService.ValidatePassword("Password-123")
	assert.NoError(t, err)

	err = authService.ValidatePassword("wrongpassword")
	assert.Error(t, err)
}
