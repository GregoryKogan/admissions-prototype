package users_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/jackc/pgx/pgtype"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	// Start a PostgreSQL container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=admissions_test",
		},
	}, func(config *docker.HostConfig) {
		// Expose the container to the host
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	require.NoError(t, err)

	// Clean up the container after tests
	t.Cleanup(func() {
		err := pool.Purge(resource)
		require.NoError(t, err)
	})

	// Exponential backoff-retry to connect to the database
	var db *gorm.DB
	err = pool.Retry(func() error {
		dsn := fmt.Sprintf("host=localhost port=%s user=postgres password=postgres dbname=admissions_test sslmode=disable",
			resource.GetPort("5432/tcp"))
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		return db.
			Exec("SELECT 1").
			Error
	})
	require.NoError(t, err)

	// Migrate the schema
	err = db.AutoMigrate(&users.User{}, &users.Role{})
	require.NoError(t, err)

	return db
}

func TestCreateRole(t *testing.T) {
	db := setupTestDB(t)
	repo := users.NewUsersRepo(db)

	role := &users.Role{
		Title:       "test_role",
		Permissions: pgtype.JSONB{Bytes: []byte(`{"read": true}`), Status: pgtype.Present},
	}
	err := repo.CreateRole(role)
	assert.NoError(t, err)

	var result users.Role
	err = db.First(&result, "title = ?", "test_role").Error
	assert.NoError(t, err)
	assert.Equal(t, "test_role", result.Title)
}

func TestRoleExists(t *testing.T) {
	db := setupTestDB(t)
	repo := users.NewUsersRepo(db)

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
	db := setupTestDB(t)
	repo := users.NewUsersRepo(db)

	role := &users.Role{Title: "test_role"}
	err := repo.CreateRole(role)
	assert.NoError(t, err)

	result, err := repo.GetRoleByTitle("test_role")
	assert.NoError(t, err)
	assert.Equal(t, "test_role", result.Title)
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := users.NewUsersRepo(db)

	repo.CreateRole(&users.Role{Title: "test_role"})
	role, _ := repo.GetRoleByTitle("test_role")

	user := &users.User{Email: "test@example.com", RoleID: role.ID, Role: *role}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	var result users.User
	err = db.First(&result, "email = ?", "test@example.com").Error
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t)
	repo := users.NewUsersRepo(db)

	repo.CreateRole(&users.Role{Title: "test_role"})
	role, _ := repo.GetRoleByTitle("test_role")

	user := &users.User{Email: "test@example.com", RoleID: role.ID, Role: *role}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	err = repo.DeleteUser(user.ID)
	assert.NoError(t, err)

	var result users.User
	err = db.First(&result, user.ID).Error
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := users.NewUsersRepo(db)

	repo.CreateRole(&users.Role{Title: "test_role"})
	role, _ := repo.GetRoleByTitle("test_role")

	user := &users.User{Email: "test@example.com", RoleID: role.ID, Role: *role}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	result, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", result.Email)
}

func TestGetWithDetailsByID(t *testing.T) {
	db := setupTestDB(t)
	repo := users.NewUsersRepo(db)

	role := &users.Role{Title: "test_role"}
	err := repo.CreateRole(role)
	assert.NoError(t, err)

	user := &users.User{Email: "test@example.com", RoleID: role.ID}
	err = repo.CreateUser(user)
	assert.NoError(t, err)

	result, err := repo.GetWithDetailsByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "test_role", result.Role.Title)
}

func TestUserExistsByID(t *testing.T) {
	db := setupTestDB(t)
	repo := users.NewUsersRepo(db)

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
	db := setupTestDB(t)
	repo := users.NewUsersRepo(db)

	repo.CreateRole(&users.Role{Title: "test_role"})
	role, _ := repo.GetRoleByTitle("test_role")

	user := &users.User{Email: "test@example.com", RoleID: role.ID, Role: *role}
	err := repo.CreateUser(user)
	assert.NoError(t, err)

	result, err := repo.GetByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", result.Email)
}
