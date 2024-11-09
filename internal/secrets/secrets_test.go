package secrets_test

import (
	"testing"

	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/stretchr/testify/assert"
)

func TestReadSecretMock(t *testing.T) {
	secretName := "test_secret"
	expectedValue := "secret_value"

	secrets.SetMockSecret(secretName, expectedValue)
	defer secrets.ClearMockSecrets()

	value, err := secrets.ReadSecret(secretName)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
}

func TestReadSecretFileNotFound(t *testing.T) {
	secretName := "non_existent_secret"

	value, err := secrets.ReadSecret(secretName)
	assert.Error(t, err)
	assert.Empty(t, value)
}

func TestSetAndClearMockSecrets(t *testing.T) {
	secretName := "another_test_secret"
	expectedValue := "another_secret_value"

	secrets.SetMockSecret(secretName, expectedValue)
	value, err := secrets.ReadSecret(secretName)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)

	secrets.ClearMockSecrets()
	value, err = secrets.ReadSecret(secretName)
	assert.Error(t, err)
	assert.Empty(t, value)
}
