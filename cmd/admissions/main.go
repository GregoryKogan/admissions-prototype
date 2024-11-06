package main

import (
	"fmt"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/config"
	"github.com/L2SH-Dev/admissions/internal/logging"
	"github.com/L2SH-Dev/admissions/internal/ping"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	config.Init()
	logging.Init()

	e := echo.New()

	addMiddleware(e)

	serveFrontend(e)

	pingHandler := ping.NewPingHandler()

	api := e.Group("/api")
	pingHandler.AddRoutes(api)

	startServer(e)
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
