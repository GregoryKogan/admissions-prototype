package main

import (
	"github.com/L2SH-Dev/admissions/internal/auth"
	"github.com/L2SH-Dev/admissions/internal/config"
	"github.com/L2SH-Dev/admissions/internal/logging"
	"github.com/L2SH-Dev/admissions/internal/passwords"
	"github.com/L2SH-Dev/admissions/internal/ping"
	"github.com/L2SH-Dev/admissions/internal/server"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/labstack/echo/v4"
)

func main() {
	config.Init()
	logging.Init()

	srv := server.NewServer()

	pingHandler := ping.NewPingHandler()
	usersHandler := users.NewUsersHandler(
		users.NewUsersService(
			users.NewUsersRepo(srv.GetDB()),
		),
		auth.NewAuthService(
			passwords.NewPasswordsService(
				passwords.NewPasswordsRepo(srv.GetDB()),
			),
		),
	)

	e := srv.GetEcho()

	serveFrontend(e)

	api := e.Group("/api")
	pingHandler.AddRoutes(api)
	usersHandler.AddRoutes(api)

	srv.Start()
}

func serveFrontend(e *echo.Echo) {
	e.Static("/", "ui/dist")
	e.File("/", "ui/dist/index.html")
}
