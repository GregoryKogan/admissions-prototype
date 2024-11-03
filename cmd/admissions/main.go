package main

import (
	"log"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/config"
	"github.com/L2SH-Dev/admissions/internal/database"
	"github.com/L2SH-Dev/admissions/internal/ping"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.Init()

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = db.AutoMigrate()
	if err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
	}

	e := echo.New()

	e.Static("/", "ui/dist")
	e.File("/", "ui/dist/index.html")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8888"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	api := e.Group("/api")
	api.GET("/ping", ping.PingHandler)

	e.Logger.Fatal(e.Start(":8888"))
}
