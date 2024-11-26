package roles_test

import (
	"testing"

	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/stretchr/testify/assert"
)

func setupTestService(t *testing.T) roles.RolesService {
	t.Cleanup(func() {
		err := storage.Flush()
		assert.NoError(t, err)
	})

	repo := roles.NewRolesRepo(storage)
	return roles.NewRolesService(repo)
}

func TestCreateRoleService(t *testing.T) {
	service := setupTestService(t)

	role := &roles.Role{Title: "test_role"}
	err := service.CreateRole(role)
	assert.NoError(t, err)

	exists, err := service.RoleExists("test_role")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestCreateDefaultRolesService(t *testing.T) {
	service := setupTestService(t)

	err := service.CreateDefaultRoles()
	assert.NoError(t, err)

	adminExists, err := service.RoleExists("admin")
	assert.NoError(t, err)
	assert.True(t, adminExists)

	userExists, err := service.RoleExists("user")
	assert.NoError(t, err)
	assert.True(t, userExists)

	interviewerExists, err := service.RoleExists("interviewer")
	assert.NoError(t, err)
	assert.True(t, interviewerExists)

	principalExists, err := service.RoleExists("principal")
	assert.NoError(t, err)
	assert.True(t, principalExists)
}

func TestRoleExistsService(t *testing.T) {
	service := setupTestService(t)

	role := &roles.Role{Title: "test_role"}
	err := service.CreateRole(role)
	assert.NoError(t, err)

	exists, err := service.RoleExists("test_role")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = service.RoleExists("non_existent_role")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestGetRoleByTitleService(t *testing.T) {
	service := setupTestService(t)

	role := &roles.Role{Title: "test_role"}
	err := service.CreateRole(role)
	assert.NoError(t, err)

	result, err := service.GetRoleByTitle("test_role")
	assert.NoError(t, err)
	assert.Equal(t, "test_role", result.Title)
}
