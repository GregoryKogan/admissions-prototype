package passwords_test

import (
	"context"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestService(t *testing.T) passwords.PasswordsService {
	repo := setupTestRepo(t)

	t.Cleanup(func() {
		err := storage.DB.Exec("DELETE FROM passwords").Error
		assert.NoError(t, err)

		err = storage.Cache.FlushDB(context.Background()).Err()
		assert.NoError(t, err)
	})

	return passwords.NewPasswordsService(repo)
}

func TestPasswordsService_Create(t *testing.T) {
	service := setupTestService(t)

	err := service.Create(1, "Valid1Password!")
	require.NoError(t, err)

	record, err := service.GetByUserID(1)
	require.NoError(t, err)
	assert.NotNil(t, record)
}

func TestPasswordsService_Validate(t *testing.T) {
	service := setupTestService(t)

	tests := []struct {
		password string
		err      error
	}{
		{"short", passwords.ErrPasswordTooShort},
		{"NoNumber!", passwords.ErrPasswordNoNumber},
		{"nonumber1", passwords.ErrPasswordNoUpper},
		{"NOLOWER1!", passwords.ErrPasswordNoLower},
		{"NoSpecial1", passwords.ErrPasswordNoSpecial},
		{"Valid1Password!", nil},
	}

	for _, tt := range tests {
		err := service.Validate(tt.password)
		if tt.err != nil {
			assert.ErrorIs(t, err, tt.err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestPasswordsService_Verify(t *testing.T) {
	service := setupTestService(t)

	err := service.Create(1, "Valid1Password!")
	require.NoError(t, err)

	match, err := service.Verify(1, "Valid1Password!")
	require.NoError(t, err)
	assert.True(t, match)

	match, err = service.Verify(1, "InvalidPassword!")
	require.NoError(t, err)
	assert.False(t, match)
}
