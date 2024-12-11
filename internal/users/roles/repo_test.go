package roles_test

import (
	"os"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	storage datastore.MockStorage
)

func TestMain(m *testing.M) {
	s, cleanup := datastore.InitMockStorage()
	storage = s

	viper.Set("users.default_role", "user")
	viper.Set("users.roles", []string{"admin", "user", "interviewer", "principal"})
	viper.Set("users.roles.user.permissions.admin", false)
	viper.Set("users.roles.user.permissions.write_general", false)
	viper.Set("users.roles.user.permissions.ai_access", false)
	viper.Set("users.roles.admin.permissions.admin", true)
	viper.Set("users.roles.admin.permissions.write_general", true)
	viper.Set("users.roles.admin.permissions.ai_access", false)
	viper.Set("users.roles.interviewer.permissions.admin", true)
	viper.Set("users.roles.interviewer.permissions.write_general", false)
	viper.Set("users.roles.interviewer.permissions.ai_access", true)
	viper.Set("users.roles.principal.permissions.admin", true)
	viper.Set("users.roles.principal.permissions.write_general", true)
	viper.Set("users.roles.principal.permissions.ai_access", true)

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func setupTestRepo(t *testing.T) roles.RolesRepo {
	t.Cleanup(func() {
		err := storage.Flush()
		assert.NoError(t, err)
	})

	return roles.NewRolesRepo(storage)
}

func TestCreateRole(t *testing.T) {
	repo := setupTestRepo(t)

	role := &roles.Role{Title: "test_role"}
	err := repo.CreateRole(role)
	assert.NoError(t, err)

	exists, err := repo.RoleExists("test_role")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestRoleExists(t *testing.T) {
	repo := setupTestRepo(t)

	role := &roles.Role{Title: "test_role"}
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

	role := &roles.Role{Title: "test_role"}
	err := repo.CreateRole(role)
	assert.NoError(t, err)

	result, err := repo.GetRoleByTitle("test_role")
	assert.NoError(t, err)
	assert.Equal(t, "test_role", result.Title)
}
