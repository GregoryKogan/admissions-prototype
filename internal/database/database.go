package database

import (
	"fmt"

	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDBConnection() *gorm.DB {
	dbConfig := viper.Sub("database")
	db_password, err := secrets.ReadSecret("db_password")
	if err != nil {
		slog.Error("Failed to connect to the database", slog.Any("error", err))
		panic(err)
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
		)), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

	if err != nil {
		slog.Error("Failed to connect to the database", slog.Any("error", err))
		panic(err)
	}

	slog.Info("Connected to the database")
	return db
}
