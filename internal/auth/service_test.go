package auth_test

import (
	"context"
	"os"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/auth"
	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/passwords"
	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/stretchr/testify/assert"
)

var (
	storage datastore.Storage
)

func TestMain(m *testing.M) {
	s, cleanup := datastore.SetupMockStorage()
	storage = s

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func setupTestService(t *testing.T) auth.AuthService {
	passwordsRepo := passwords.NewPasswordsRepo(storage)
	passwordsService := passwords.NewPasswordsService(passwordsRepo)

	t.Cleanup(func() {
		err := storage.DB.Exec("DELETE FROM passwords").Error
		assert.NoError(t, err)

		err = storage.Cache.FlushDB(context.Background()).Err()
		assert.NoError(t, err)
	})

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
