package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/L2SH-Dev/admissions/internal/config"
	"github.com/L2SH-Dev/admissions/internal/database"
	"github.com/L2SH-Dev/admissions/internal/logging"
	"github.com/L2SH-Dev/admissions/internal/ping"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func main() {
	config.Init()
	logging.Init()

	initDatabase()

	e := echo.New()

	addMiddleware(e)

	e.Static("/", "ui/dist")
	e.File("/", "ui/dist/index.html")

	api := e.Group("/api")
	ping.AddRoutes(api)

	port := viper.GetString("server.port")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func initDatabase() {
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
}

func addMiddleware(e *echo.Echo) {
	logging.AddMiddleware(e)

	port := viper.GetString("server.port")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{fmt.Sprintf("http://localhost:%s", port)},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
}
