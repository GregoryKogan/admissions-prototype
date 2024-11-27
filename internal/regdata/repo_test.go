package regdata_test

import (
	"os"
	"testing"
	"time"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/regdata"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	storage datastore.MockStorage
)

func TestMain(m *testing.M) {
	viper.Set("secrets.jwt_key", "test_key")

	s, cleanup := datastore.InitMockStorage()
	storage = s

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func setupTestRepo(t *testing.T) regdata.RegistrationDataRepo {
	t.Cleanup(func() {
		err := storage.Flush()
		assert.NoError(t, err)
	})

	return regdata.NewRegistrationDataRepo(storage)
}

func TestCreateRegistrationData(t *testing.T) {
	repo := setupTestRepo(t)

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

	err := repo.CreateRegistrationData(data)
	assert.NoError(t, err)
	assert.NotZero(t, data.ID)

	// Verify data was created
	result, err := repo.GetByID(data.ID)
	assert.NoError(t, err)
	assert.Equal(t, data.Email, result.Email)
	assert.Equal(t, data.FirstName, result.FirstName)
	assert.Equal(t, data.Grade, result.Grade)
}

func TestGetByID(t *testing.T) {
	repo := setupTestRepo(t)

	// Test getting non-existent record
	_, err := repo.GetByID(999)
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
	err = repo.CreateRegistrationData(data)
	require.NoError(t, err)

	// Test getting existing record
	result, err := repo.GetByID(data.ID)
	assert.NoError(t, err)
	assert.Equal(t, data.Email, result.Email)
	assert.Equal(t, data.FirstName, result.FirstName)
	assert.Equal(t, data.Grade, result.Grade)
}

func TestExistsByEmailNameAndGrade(t *testing.T) {
	repo := setupTestRepo(t)

	// Test non-existent record
	exists, err := repo.ExistsByEmailNameAndGrade("nonexistent@example.com", "Test", 9)
	assert.NoError(t, err)
	assert.False(t, exists)

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
	err = repo.CreateRegistrationData(data)
	require.NoError(t, err)

	// Test existing record
	exists, err = repo.ExistsByEmailNameAndGrade(data.Email, data.FirstName, data.Grade)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestSetEmailVerified(t *testing.T) {
	repo := setupTestRepo(t)

	// Test non-existent record
	err := repo.SetEmailVerified(999)
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
	err = repo.CreateRegistrationData(data)
	require.NoError(t, err)

	// Test setting email verified
	err = repo.SetEmailVerified(data.ID)
	assert.NoError(t, err)

	// Verify email was marked as verified
	result, err := repo.GetByID(data.ID)
	assert.NoError(t, err)
	assert.True(t, result.EmailVerified)
}