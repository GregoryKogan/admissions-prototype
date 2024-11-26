package roles_test

import (
	"os"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/jackc/pgx/pgtype"
	"github.com/stretchr/testify/assert"
)

var (
	storage datastore.Storage
)

func TestMain(m *testing.M) {
	s, cleanup := datastore.InitMockStorage()
	storage = s

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func setupTestRepo(t *testing.T) roles.RolesRepo {
	t.Cleanup(func() {
		err := storage.DB().Exec("DELETE FROM roles").Error
		assert.NoError(t, err)
	})

	return roles.NewRolesRepo(storage)
}

func TestCreateRole(t *testing.T) {
	repo := setupTestRepo(t)

	role := &roles.Role{
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
