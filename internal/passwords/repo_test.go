package passwords_test

import (
	"testing"

	"github.com/L2SH-Dev/admissions/internal/passwords"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	return db
}

func TestNewPasswordsRepo(t *testing.T) {
	db := setupTestDB(t)
	repo := passwords.NewPasswordsRepo(db)
	assert.NotNil(t, repo)
}

func TestCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := passwords.NewPasswordsRepo(db)

	hashedPassword := &passwords.HashedPassword{
		Hash:      []byte("hashedpassword"),
		Salt:      []byte("salt"),
		Algorithm: "argon2id",
	}

	err := repo.Create(1, hashedPassword)
	require.NoError(t, err)

	var record passwords.Password
	err = db.Where("user_id = ?", 1).First(&record).Error
	require.NoError(t, err)
	assert.Equal(t, hashedPassword.Hash, record.Hash)
	assert.Equal(t, hashedPassword.Salt, record.Salt)
	assert.Equal(t, hashedPassword.Algorithm, record.Algorithm)
}

func TestGetByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := passwords.NewPasswordsRepo(db)

	hashedPassword := &passwords.HashedPassword{
		Hash:      []byte("hashedpassword"),
		Salt:      []byte("salt"),
		Algorithm: "argon2id",
	}

	err := repo.Create(1, hashedPassword)
	require.NoError(t, err)

	record, err := repo.GetByUserID(1)
	require.NoError(t, err)
	assert.Equal(t, hashedPassword.Hash, record.Hash)
	assert.Equal(t, hashedPassword.Salt, record.Salt)
	assert.Equal(t, hashedPassword.Algorithm, record.Algorithm)
}

func TestExistsByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := passwords.NewPasswordsRepo(db)

	hashedPassword := &passwords.HashedPassword{
		Hash:      []byte("hashedpassword"),
		Salt:      []byte("salt"),
		Algorithm: "argon2id",
	}

	err := repo.Create(1, hashedPassword)
	require.NoError(t, err)

	exists, err := repo.ExistsByUserID(1)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByUserID(2)
	require.NoError(t, err)
	assert.False(t, exists)
}
