package users_test

import (
	"os"
	"testing"
	"time"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/regdata"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
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
	viper.Set("secrets.jwt_key", "test_key")
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

func setupTestRepo(t *testing.T) users.UsersRepo {
	t.Cleanup(func() {
		err := storage.Flush()
		assert.NoError(t, err)
	})

	roles.NewRolesService(roles.NewRolesRepo(storage)).CreateDefaultRoles()

	rolesRepo := roles.NewRolesRepo(storage)
	rolesService := roles.NewRolesService(rolesRepo)
	usersRepo := users.NewUsersRepo(storage)
	usersService := users.NewUsersService(usersRepo, rolesService)

	passwordsRepo := passwords.NewPasswordsRepo(storage)
	passwordsService := passwords.NewPasswordsService(passwordsRepo)
	authRepo := auth.NewAuthRepo(storage)
	authService := auth.NewAuthService(authRepo, passwordsService)

	repo := regdata.NewRegistrationDataRepo(storage)

	regdataService := regdata.NewRegistrationDataService(repo, usersService, authService, passwordsService)
	err := regdataService.Create(&regdata.RegistrationData{
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
	err = regdataService.Create(&regdata.RegistrationData{
		Email:           "test2@mail.org",
		FirstName:       "Test",
		LastName:        "User2",
		Gender:          "F",
		BirthDate:       time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
		Grade:           6,
		OldSchool:       "Test School 2",
		ParentFirstName: "Test",
		ParentLastName:  "Parent2",
		ParentPhone:     "+79999999999",
	})
	assert.NoError(t, err)
	return users.NewUsersRepo(storage)
}

func TestCreateUser(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
	err := repo.Create(user)
	assert.NoError(t, err)

	exists, err := repo.ExistsByID(user.ID)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestDeleteUser(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
	err := repo.Create(user)
	assert.NoError(t, err)

	err = repo.Delete(user.ID)
	assert.NoError(t, err)

	exists, err := repo.ExistsByID(user.ID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestGetByID(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
	err := repo.Create(user)
	assert.NoError(t, err)

	result, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Login, result.Login)
	assert.Equal(t, user.RoleID, result.RoleID)
	assert.Equal(t, user.RegistrationDataID, result.RegistrationDataID)
	assert.Equal(t, uint(1), result.RegistrationDataID)
	assert.Equal(t, uint(1), result.Role.ID)
}

func TestUserExistsByID(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
	err := repo.Create(user)
	assert.NoError(t, err)

	exists, err := repo.ExistsByID(user.ID)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByID(999)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestGetByRegistrationID(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
	err := repo.Create(user)
	assert.NoError(t, err)

	result, err := repo.GetByRegistrationID(user.RegistrationDataID)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Login, result.Login)
	assert.Equal(t, user.RoleID, result.RoleID)
	assert.Equal(t, user.RegistrationDataID, result.RegistrationDataID)
	assert.Equal(t, uint(1), result.RegistrationDataID)
	assert.Equal(t, uint(1), result.Role.ID)
}

func TestGetByLogin(t *testing.T) {
	repo := setupTestRepo(t)

	user := &users.User{Login: "test_login", RoleID: 1, RegistrationDataID: 1}
	err := repo.Create(user)
	assert.NoError(t, err)

	result, err := repo.GetByLogin(user.Login)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Login, result.Login)
	assert.Equal(t, user.RoleID, result.RoleID)
	assert.Equal(t, user.RegistrationDataID, result.RegistrationDataID)
	assert.Equal(t, uint(1), result.RegistrationDataID)
	assert.Equal(t, uint(1), result.Role.ID)
}
