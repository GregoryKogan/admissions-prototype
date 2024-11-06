package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/auth"
	"github.com/L2SH-Dev/admissions/internal/config"
	"github.com/L2SH-Dev/admissions/internal/database"
	"github.com/L2SH-Dev/admissions/internal/logging"
	"github.com/L2SH-Dev/admissions/internal/passwords"
	"github.com/L2SH-Dev/admissions/internal/ping"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	config.Init()
	logging.Init()

	db := initDB()

	e := echo.New()

	addMiddleware(e)

	validation.AddValidation(e)

	serveFrontend(e)

	pingHandler := ping.NewPingHandler()
	usersHandler := users.NewUsersHandler(
		users.NewUsersService(
			users.NewUsersRepo(db),
		),
		auth.NewAuthService(
			passwords.NewPasswordsService(
				passwords.NewPasswordsRepo(db),
			),
		),
	)

	api := e.Group("/api")
	pingHandler.AddRoutes(api)
	usersHandler.AddRoutes(api)

	startServer(e)
}

func initDB() *gorm.DB {
	db, err := database.Connect()
	if err != nil {
		slog.Error("Failed to connect to the database", slog.Any("error", err))
		panic(err)
	}

	return db
}

func addMiddleware(e *echo.Echo) {
	logging.AddMiddleware(e)

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	port := viper.GetString("server.port")
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{fmt.Sprintf("http://localhost:%s", port)},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
}

func serveFrontend(e *echo.Echo) {
	e.Static("/", "ui/dist")
	e.File("/", "ui/dist/index.html")
}

func startServer(e *echo.Echo) {
	port := viper.GetString("server.port")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
