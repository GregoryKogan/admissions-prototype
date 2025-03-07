package passwords_test

import (
	"os"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords/crypto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	storage datastore.MockStorage
)

func TestMain(m *testing.M) {
	viper.Set("auth.passwords.min_length", 8)

	s, cleanup := datastore.InitMockStorage()
	storage = s

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func setupTestRepo(t *testing.T) passwords.PasswordsRepo {
	t.Cleanup(func() {
		err := storage.Flush()
		assert.NoError(t, err)
	})

	return passwords.NewPasswordsRepo(storage)
}

func TestNewPasswordsRepo(t *testing.T) {
	repo := setupTestRepo(t)
	assert.NotNil(t, repo)
}

func TestCreate(t *testing.T) {
	repo := setupTestRepo(t)

	hashedPassword := &crypto.HashedPassword{
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

func TestGetByUserID(t *testing.T) {
	repo := setupTestRepo(t)

	hashedPassword := &crypto.HashedPassword{
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
	repo := setupTestRepo(t)

	hashedPassword := &crypto.HashedPassword{
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
