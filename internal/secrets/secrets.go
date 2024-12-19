package secrets

import (
	"fmt"
	"os"
	"strings"

	"log/slog"

	"github.com/spf13/viper"
)

func LoadSecretsIntoViper() error {
	secrets := []string{
		"db_password",
		"jwt_key",
		"mail_api_key",
		"admin_password",
	}

	for _, secretName := range secrets {
		envVarName := toEnvVarName(secretName)
		value := os.Getenv(envVarName)
		if value == "" {
			slog.Error("Environment variable not set", slog.String("variable", envVarName))
			return fmt.Errorf("environment variable %s not set", envVarName)
		}
		viper.Set("secrets."+secretName, value)
		slog.Debug("Loaded secret", slog.String("name", secretName))
	}

	slog.Info("All secrets loaded into viper from environment variables")
	return nil
}

func toEnvVarName(secretName string) string {
	// Convert secret_name to SECRET_NAME
	return strings.ToUpper(strings.ReplaceAll(secretName, ".", "_"))
}
