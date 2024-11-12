package auth_test

import (
	"context"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/users/auth"
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
