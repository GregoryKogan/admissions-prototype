package main

import (
	"net/http"
	"os"

	"github.com/L2SH-Dev/admissions/internal/config"
	"github.com/L2SH-Dev/admissions/internal/database"
	"github.com/L2SH-Dev/admissions/internal/logging"
	"github.com/L2SH-Dev/admissions/internal/ping"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/exp/slog"
)

func main() {
	config.Init()
	logging.Init()

	db, err := database.Connect()
	if err != nil {
		slog.Error("Failed to connect to the database", slog.Any("error", err))
		os.Exit(1)
	}

	err = db.AutoMigrate()
	if err != nil {
		slog.Error("Failed to migrate the database", slog.Any("error", err))
		os.Exit(1)
	}

	e := echo.New()

	logging.AddMiddleware(e)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8888"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	e.Static("/", "ui/dist")
	e.File("/", "ui/dist/index.html")

	api := e.Group("/api")
	api.GET("/ping", ping.PingHandler)

	e.Logger.Fatal(e.Start(":8888"))
}
