package config_test

import (
	"os"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitConfigFileNotFound(t *testing.T) {
	viper.Reset()
	viper.SetConfigName("nonexistent")
	viper.AddConfigPath(".")
	assert.Panics(t, func() {
		config.Init()
	})
}

func TestInitConfigFileFound(t *testing.T) {
	viper.Reset()
	viper.SetConfigName("config")
	viper.WriteConfigAs("config.yaml")
	viper.AddConfigPath(".")
	assert.NotPanics(t, func() {
		config.Init()
	})
	os.Remove("config.yaml")
}
