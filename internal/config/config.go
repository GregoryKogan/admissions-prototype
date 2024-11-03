package config

import (
	"os"

	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Error("Config file not found", slog.Any("error", err))
			os.Exit(1)
		} else {
			slog.Error("Error reading config file", slog.Any("error", err))
			os.Exit(1)
		}
	}

	slog.Info("Config file successfully read")
}
