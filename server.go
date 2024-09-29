package main

import (
	"net/http"

	"github.com/L2SH-Dev/admissions/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Static("/", "ui/dist")
	e.File("/", "ui/dist/index.html")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8888"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	api := e.Group("/api")
	api.GET("/ping", handlers.PingHandler)

	e.Logger.Fatal(e.Start(":8888"))
}
