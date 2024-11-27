package regdata_test

import (
	"testing"
	"time"

	"github.com/L2SH-Dev/admissions/internal/regdata"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestService(t *testing.T) regdata.RegistrationDataService {
	t.Cleanup(func() {
		err := storage.Flush()
		assert.NoError(t, err)
	})

	rolesRepo := roles.NewRolesRepo(storage)
	rolesService := roles.NewRolesService(rolesRepo)
	err := rolesService.CreateDefaultRoles()
	require.NoError(t, err)

	usersRepo := users.NewUsersRepo(storage)
	usersService := users.NewUsersService(usersRepo, rolesService)

	passwordsRepo := passwords.NewPasswordsRepo(storage)
	passwordsService := passwords.NewPasswordsService(passwordsRepo)
	authRepo := auth.NewAuthRepo(storage)
	authService := auth.NewAuthService(authRepo, passwordsService)

	repo := regdata.NewRegistrationDataRepo(storage)
	return regdata.NewRegistrationDataService(repo, usersService, authService, passwordsService)
}

func TestCreateRegistrationDataService(t *testing.T) {
	service := setupTestService(t)

	// Test valid registration data
	data := &regdata.RegistrationData{
		Email:           "test@example.com",
		FirstName:       "Test",
		LastName:        "User",
		Gender:          "M",
		BirthDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Grade:           9,
		OldSchool:       "Previous School",
		ParentFirstName: "Parent",
		ParentLastName:  "Test",
		ParentPhone:     "+1234567890",
	}
	err := service.CreateRegistrationData(data)
	assert.NoError(t, err)
	assert.NotZero(t, data.ID)

	// Test invalid data (missing required fields)
	invalidData := &regdata.RegistrationData{
		Email: "test@example.com",
	}
	err = service.CreateRegistrationData(invalidData)
	assert.ErrorIs(t, err, regdata.ErrRegistrationDataInvalid)

	// Test duplicate registration data
	duplicateData := &regdata.RegistrationData{
		Email:           "test@example.com",
		FirstName:       "Test",
		Grade:           9,
		LastName:        "User2",
		Gender:          "M",
		BirthDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		OldSchool:       "Previous School",
		ParentFirstName: "Parent",
		ParentLastName:  "Test",
		ParentPhone:     "+1234567890",
	}
	err = service.CreateRegistrationData(duplicateData)
	assert.ErrorIs(t, err, regdata.ErrRegistrationDataExists)
}

func TestGetByIDService(t *testing.T) {
	service := setupTestService(t)

	// Test getting non-existent record
	_, err := service.GetByID(999)
	assert.Error(t, err)

	// Create test data
	data := &regdata.RegistrationData{
		Email:           "test@example.com",
		FirstName:       "Test",
		LastName:        "User",
		Gender:          "M",
		BirthDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Grade:           9,
		OldSchool:       "Previous School",
		ParentFirstName: "Parent",
		ParentLastName:  "Test",
		ParentPhone:     "+1234567890",
	}
	err = service.CreateRegistrationData(data)
	require.NoError(t, err)

	// Test getting existing record
	result, err := service.GetByID(data.ID)
	assert.NoError(t, err)
	assert.Equal(t, data.Email, result.Email)
	assert.Equal(t, data.FirstName, result.FirstName)
}

func TestSetEmailVerifiedService(t *testing.T) {
	service := setupTestService(t)

	// Create test data
	data := &regdata.RegistrationData{
		Email:           "test@example.com",
		FirstName:       "Test",
		LastName:        "User",
		Gender:          "M",
		BirthDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Grade:           9,
		OldSchool:       "Previous School",
		ParentFirstName: "Parent",
		ParentLastName:  "Test",
		ParentPhone:     "+1234567890",
	}
	err := service.CreateRegistrationData(data)
	require.NoError(t, err)

	// Test setting email verified
	err = service.SetEmailVerified(data.ID)
	assert.NoError(t, err)

	// Verify email was marked as verified
	result, err := service.GetByID(data.ID)
	assert.NoError(t, err)
	assert.True(t, result.EmailVerified)

	// Test setting non-existent record
	err = service.SetEmailVerified(999)
	assert.Error(t, err)
}

func TestAcceptService(t *testing.T) {
	service := setupTestService(t)

	// Test accepting non-existent registration
	_, err := service.Accept(999)
	assert.Error(t, err)

	// Create test data
	data := &regdata.RegistrationData{
		Email:           "test@example.com",
		FirstName:       "Test",
		LastName:        "User",
		Gender:          "M",
		BirthDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Grade:           9,
		OldSchool:       "Previous School",
		ParentFirstName: "Parent",
		ParentLastName:  "Test",
		ParentPhone:     "+1234567890",
	}
	err = service.CreateRegistrationData(data)
	require.NoError(t, err)

	// Test accepting unverified email
	_, err = service.Accept(data.ID)
	assert.ErrorIs(t, err, regdata.ErrorEmailNotVerified)

	// Verify email and test successful acceptance
	err = service.SetEmailVerified(data.ID)
	require.NoError(t, err)

	user, err := service.Accept(data.ID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, data.ID, user.RegistrationDataID)
	assert.Equal(t, "Test User", user.Login)
}
