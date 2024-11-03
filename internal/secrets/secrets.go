package secrets

import "os"

func ReadSecret(secret string) (string, error) {
	buffer, err := os.ReadFile("/run/secrets/" + secret)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}
