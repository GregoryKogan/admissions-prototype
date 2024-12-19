package users_test

import (
	"testing"

	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestService(t *testing.T) users.UsersService {
	t.Cleanup(func() {
		err := storage.Flush()
		assert.NoError(t, err)
	})

	repo := setupTestRepo(t)
	rolesService := roles.NewRolesService(roles.NewRolesRepo(storage))
	return users.NewUsersService(repo, rolesService)
}

func TestUsersService_CreateUser(t *testing.T) {
	service := setupTestService(t)

	createdUser, err := service.Create(1, "test_login")
	assert.NoError(t, err)
	assert.Equal(t, "test_login", createdUser.Login)

	// Check if user was created
	user, err := service.GetByLogin("test_login")
	assert.NoError(t, err)
	assert.Equal(t, "test_login", user.Login)

	// Check if user has default role
	assert.Equal(t, "user", user.Role.Title)

	// Create second user
	_, err = service.Create(2, "test2_login")
	assert.NoError(t, err)
}

func TestUsersService_CreateUserLoginExists(t *testing.T) {
	service := setupTestService(t)

	_, err := service.Create(1, "test_login")
	assert.NoError(t, err)

	_, err = service.Create(2, "test_login")
	assert.Error(t, err)
}

func TestUsersService_CreateUserRegistrationDataExists(t *testing.T) {
	service := setupTestService(t)

	_, err := service.Create(1, "test_login")
	assert.NoError(t, err)

	_, err = service.Create(1, "test2_login")
	assert.Error(t, err)
}

func TestUsersService_DeleteUser(t *testing.T) {
	service := setupTestService(t)

	createdUser, err := service.Create(1, "test_login")
	assert.NoError(t, err)

	err = service.Delete(createdUser.ID)
	assert.NoError(t, err)

	_, err = service.GetByLogin("test_login")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUsersService_GetByID(t *testing.T) {
	service := setupTestService(t)

	createdUser, err := service.Create(1, "test_login")
	assert.NoError(t, err)

	user, err := service.GetByID(createdUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, "test_login", user.Login)
}

func TestUsersService_GetByIDNotFound(t *testing.T) {
	service := setupTestService(t)

	_, err := service.Create(1, "test_login")
	assert.NoError(t, err)

	_, err = service.GetByID(999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUsersService_GetByLogin(t *testing.T) {
	service := setupTestService(t)

	createdUser, err := service.Create(1, "test_login")
	assert.NoError(t, err)

	user, err := service.GetByLogin("test_login")
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, "test_login", user.Login)
}

func TestUsersService_GetByLoginNotFound(t *testing.T) {
	service := setupTestService(t)

	_, err := service.Create(1, "test_login")
	assert.NoError(t, err)

	_, err = service.GetByLogin("not_found")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
