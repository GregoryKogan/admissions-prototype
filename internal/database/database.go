package database

import (
	"fmt"

	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dbConfig := viper.Sub("database")
	db_password, err := secrets.ReadSecret("db_password")
	if err != nil {
		slog.Error("Failed to read database password", slog.Any("error", err))
		return nil, err
	}

	db, err := gorm.Open(
		postgres.Open(fmt.Sprintf(
			"host=%s port=%v user=%s dbname=%s password=%s sslmode=disable TimeZone=%s",
			dbConfig.GetString("host"),
			dbConfig.GetInt("port"),
			dbConfig.GetString("user"),
			dbConfig.GetString("name"),
			db_password,
			dbConfig.GetString("connection.timezone"),
		)), &gorm.Config{})

	if err != nil {
		slog.Error("Failed to connect to the database", slog.Any("error", err))
		return nil, err
	}

	slog.Info("Connected to the database")
	return db, nil
}
