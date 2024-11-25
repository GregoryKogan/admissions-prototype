package users_test

import (
	"context"
	"os"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	storage datastore.Storage
)

func TestMain(m *testing.M) {
	s, cleanup := datastore.SetupMockStorage()
	storage = s

	viper.Set("auth.access_lifetime", "15m")
	viper.Set("auth.refresh_lifetime", "720h")
	viper.Set("auth.auto_logout", "24h")

	secrets.SetMockSecret("jwt_key", "test_key")

	roles.NewRolesService(roles.NewRolesRepo(storage)).CreateDefaultRoles()

	code := m.Run()

	secrets.ClearMockSecrets()
	if err := storage.DB.Exec("DELETE FROM roles").Error; err != nil {
		panic(err)
	}
	cleanup()

	os.Exit(code)
}

func setupTestRepo(t *testing.T) users.UsersRepo {
	t.Cleanup(func() {
		err := storage.DB.Exec("DELETE FROM users").Error
		assert.NoError(t, err)

		err = storage.Cache.FlushDB(context.Background()).Err()
		assert.NoError(t, err)
	})

	return users.NewUsersRepo(storage)
}

func TestCreateUser(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Email: "test@example.com", RoleID: 1}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	exists, err := repo.UserExistsByID(user.ID)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestDeleteUser(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Email: "test@example.com", RoleID: 1}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	err = repo.DeleteUser(user.ID)
	assert.NoError(t, err)

	exists, err := repo.UserExistsByID(user.ID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestGetByID(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Email: "test@example.com", RoleID: 1}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	result, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestUserExistsByID(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Email: "test@example.com", RoleID: 1}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	exists, err := repo.UserExistsByID(user.ID)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.UserExistsByID(999)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestGetByEmail(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Email: "test@example.com", RoleID: 1}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	result, err := repo.GetByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", result.Email)
}
