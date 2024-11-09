package users_test

import (
	"errors"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUsersService_CreateUser(t *testing.T) {
	db := setupTestDB(t) // Reuse setupTestDB from repo_test.go
	repo := users.NewUsersRepo(db)
	service := users.NewUsersService(repo)

	// Create default roles
	createdUser, err := service.Create("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", createdUser.Email)

	// Check if user was created
	user, err := service.GetByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUsersService_CreateUserAlreadyExists(t *testing.T) {
	db := setupTestDB(t) // Reuse setupTestDB from repo_test.go
	repo := users.NewUsersRepo(db)
	service := users.NewUsersService(repo)

	// Create default roles
	_, err := service.Create("test@example.com")
	assert.NoError(t, err)

	// Try to create the same user again
	_, err = service.Create("test@example.com")
	assert.Error(t, err)
	assert.Equal(t, users.ErrUserAlreadyExists, err)
}

func TestUsersService_DeleteUser(t *testing.T) {
	db := setupTestDB(t) // Reuse setupTestDB from repo_test.go
	repo := users.NewUsersRepo(db)
	service := users.NewUsersService(repo)

	// Create user
	user, err := service.Create("test@example.com")
	assert.NoError(t, err)

	// Delete user
	err = service.Delete(user.ID)
	assert.NoError(t, err)

	// Check if user was deleted
	_, err = service.GetByID(user.ID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestUsersService_GetByEmail(t *testing.T) {
	db := setupTestDB(t) // Reuse setupTestDB from repo_test.go
	repo := users.NewUsersRepo(db)
	service := users.NewUsersService(repo)

	// Create user
	_, err := service.Create("test@example.com")
	assert.NoError(t, err)

	// Get user by email
	user, err := service.GetByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUsersService_GetByID(t *testing.T) {
	db := setupTestDB(t) // Reuse setupTestDB from repo_test.go
	repo := users.NewUsersRepo(db)
	service := users.NewUsersService(repo)

	// Create user
	createdUser, err := service.Create("test@example.com")
	assert.NoError(t, err)

	// Get user by ID
	user, err := service.GetByID(createdUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
}

func TestUsersService_GetFullByID(t *testing.T) {
	db := setupTestDB(t) // Reuse setupTestDB from repo_test.go
	repo := users.NewUsersRepo(db)
	service := users.NewUsersService(repo)

	// Create user
	createdUser, err := service.Create("test@example.com")
	assert.NoError(t, err)

	// Get user with details by ID
	user, err := service.GetFullByID(createdUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "user", user.Role.Title)
}
