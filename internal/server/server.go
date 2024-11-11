package server

import (
	"fmt"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/logging"
	"github.com/L2SH-Dev/admissions/internal/storage"
	"github.com/L2SH-Dev/admissions/internal/validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

type Server interface {
	Start()
	AddFrontend(static string, index string)
	AddHandlers(storage storage.Storage, handlers ...func(storage storage.Storage) Handler)
}

type server struct {
	Echo *echo.Echo
}

type Handler interface {
	AddRoutes(g *echo.Group)
}

func NewServer() Server {
	srv := &server{
		Echo: echo.New(),
	}

	srv.addMiddleware()
	validation.AddValidation(srv.Echo)

	return srv
}

func (s *server) Start() {
	port := viper.GetString("server.port")
	s.Echo.Logger.Fatal(s.Echo.Start(fmt.Sprintf(":%s", port)))
}

func (s *server) AddFrontend(static string, index string) {
	s.Echo.Static("/", static)
	s.Echo.File("/", index)
}

func (s *server) AddHandlers(storage storage.Storage, handlers ...func(storage storage.Storage) Handler) {
	for _, handler := range handlers {
		s.addHandler(handler(storage))
	}
}

func (s *server) addHandler(h Handler) {
	h.AddRoutes(s.Echo.Group("/api"))
}

func (s *server) addMiddleware() {
	logging.AddMiddleware(s.Echo)

	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.Secure())

	port := viper.GetString("server.port")
	s.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{fmt.Sprintf("http://localhost:%s", port)},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
}
