package users_test

import (
	"os"
	"testing"
	"time"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/regdata"
	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/L2SH-Dev/admissions/internal/users"
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
		err := storage.Flush()
		assert.NoError(t, err)
	})

	roles.NewRolesService(roles.NewRolesRepo(storage)).CreateDefaultRoles()
	regdataService := regdata.NewRegistrationDataService(regdata.NewRegistrationDataRepo(storage))
	err := regdataService.CreateRegistrationData(&regdata.RegistrationData{
		Email:           "test@mail.org",
		FirstName:       "Test",
		LastName:        "User",
		Gender:          "M",
		BirthDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Grade:           8,
		OldSchool:       "Test School",
		ParentFirstName: "Test",
		ParentLastName:  "Parent",
		ParentPhone:     "+79999999999",
	})
	assert.NoError(t, err)
	return users.NewUsersRepo(storage)
}

func TestCreateUser(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	exists, err := repo.UserExistsByID(user.ID)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestDeleteUser(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
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

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	result, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Login, result.Login)
	assert.Equal(t, user.RoleID, result.RoleID)
	assert.Equal(t, user.RegistrationDataID, result.RegistrationDataID)
}

func TestUserExistsByID(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	exists, err := repo.UserExistsByID(user.ID)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.UserExistsByID(999)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestGetByRegistrationID(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	result, err := repo.GetByRegistrationID(user.RegistrationDataID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Login, result.Login)
	assert.Equal(t, user.RoleID, result.RoleID)
	assert.Equal(t, user.RegistrationDataID, result.RegistrationDataID)
}

func TestGetByLogin(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	result, err := repo.GetByLogin(user.Login)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Login, result.Login)
	assert.Equal(t, user.RoleID, result.RoleID)
	assert.Equal(t, user.RegistrationDataID, result.RegistrationDataID)
}
