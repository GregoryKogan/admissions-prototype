package server

import (
	"fmt"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/logging"
	"github.com/L2SH-Dev/admissions/internal/validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

type Server interface {
	Start()
	AddFrontend(static string)
	AddHandlers(
		storage datastore.Storage,
		handlers ...func(storage datastore.Storage) Handler,
	)
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

	srv.addGeneralMiddleware()
	validation.AddValidation(srv.Echo)

	return srv
}

func (s *server) Start() {
	port := viper.GetString("server.port")
	s.Echo.Logger.Fatal(s.Echo.Start(fmt.Sprintf(":%s", port)))
}

func (s *server) AddFrontend(static string) {
	s.Echo.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  static,
		HTML5: true,
	}))
}

func (s *server) AddHandlers(
	storage datastore.Storage,
	handlers ...func(storage datastore.Storage) Handler,
) {
	for _, handler := range handlers {
		s.addHandler(handler(storage))
	}
}

func (s *server) addHandler(h Handler) {
	h.AddRoutes(s.Echo.Group("/api"))
}

func (s *server) addGeneralMiddleware() {
	logging.AddMiddleware(s.Echo)

	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.Secure())

	// Configure CORS to allow frontend requests
	s.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000", // Dev frontend
			fmt.Sprintf("http://localhost:%s", viper.GetString("server.port")), // Prod frontend
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true,
	}))
}
