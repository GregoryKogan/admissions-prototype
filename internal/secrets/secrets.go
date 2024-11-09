package secrets

import (
	"os"
)

// Mock secrets for testing
var mockSecrets = make(map[string]string)

func ReadSecret(secret string) (string, error) {
	if value, ok := mockSecrets[secret]; ok {
		return value, nil
	}

	buffer, err := os.ReadFile("/run/secrets/" + secret)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

func SetMockSecret(secret, value string) {
	mockSecrets[secret] = value
}

func ClearMockSecrets() {
	mockSecrets = make(map[string]string)
}
