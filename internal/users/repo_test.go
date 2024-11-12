package users_test

import (
	"context"
	"os"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/jackc/pgx/pgtype"
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

	code := m.Run()

	secrets.ClearMockSecrets()
	cleanup()
	os.Exit(code)
}

func setupTestRepo(t *testing.T) users.UsersRepo {
	t.Cleanup(func() {
		err := storage.DB.Exec("DELETE FROM users").Error
		assert.NoError(t, err)

		err = storage.DB.Exec("DELETE FROM roles").Error
		assert.NoError(t, err)

		err = storage.DB.Exec("DELETE FROM passwords").Error
		assert.NoError(t, err)

		err = storage.Cache.FlushDB(context.Background()).Err()
		assert.NoError(t, err)
	})

	return users.NewUsersRepo(storage)
}

func TestCreateRole(t *testing.T) {
	repo := setupTestRepo(t)

	role := &users.Role{
		Title:       "test_role",
		Permissions: pgtype.JSONB{Bytes: []byte(`{"read": true}`), Status: pgtype.Present},
	}
	err := repo.CreateRole(role)
	assert.NoError(t, err)

	exists, err := repo.RoleExists("test_role")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRoleExists(t *testing.T) {
	repo := setupTestRepo(t)

	role := &users.Role{Title: "test_role"}
	err := repo.CreateRole(role)
	assert.NoError(t, err)

	exists, err := repo.RoleExists("test_role")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.RoleExists("non_existent_role")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestGetRoleByTitle(t *testing.T) {
	repo := setupTestRepo(t)

	role := &users.Role{Title: "test_role"}
	err := repo.CreateRole(role)
	assert.NoError(t, err)

	result, err := repo.GetRoleByTitle("test_role")
	assert.NoError(t, err)
	assert.Equal(t, "test_role", result.Title)
}

func TestCreateUser(t *testing.T) {
	repo := setupTestRepo(t)

	repo.CreateRole(&users.Role{Title: "test_role"})
	role, _ := repo.GetRoleByTitle("test_role")

	user := &users.User{Email: "test@example.com", RoleID: role.ID, Role: *role}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	exists, err := repo.UserExistsByID(user.ID)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestDeleteUser(t *testing.T) {
	repo := setupTestRepo(t)

	repo.CreateRole(&users.Role{Title: "test_role"})
	role, _ := repo.GetRoleByTitle("test_role")

	user := &users.User{Email: "test@example.com", RoleID: role.ID, Role: *role}
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

	role := &users.Role{Title: "test_role"}
	err := repo.CreateRole(role)
	assert.NoError(t, err)

	user := &users.User{Email: "test@example.com", RoleID: role.ID}
	err = repo.CreateUser(user)
	assert.NoError(t, err)

	result, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "test_role", result.Role.Title)
}

func TestUserExistsByID(t *testing.T) {
	repo := setupTestRepo(t)

	repo.CreateRole(&users.Role{Title: "test_role"})
	role, _ := repo.GetRoleByTitle("test_role")

	user := &users.User{Email: "test@example.com", RoleID: role.ID, Role: *role}
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

	repo.CreateRole(&users.Role{Title: "test_role"})
	role, _ := repo.GetRoleByTitle("test_role")

	user := &users.User{Email: "test@example.com", RoleID: role.ID, Role: *role}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	result, err := repo.GetByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", result.Email)
}
