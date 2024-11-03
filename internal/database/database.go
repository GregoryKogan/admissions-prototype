package database

import (
	"fmt"
	"log"

	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dbConfig := viper.Sub("database")
	db_password, err := secrets.ReadSecret("db_password")
	if err != nil {
		log.Printf("Failed to read database password: %v", err)
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
		log.Printf("Failed to connect to the database: %v", err)
		return nil, err
	}

	return db, nil
}
