package server

import (
	"fmt"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/logging"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/L2SH-Dev/admissions/internal/validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

type AdminMiddlewareService interface {
	Add(g *echo.Group, minimalRole roles.Role) error
}

type Server interface {
	Start()
	AddFrontend(static string, index string)
	AddHandlers(
		storage datastore.Storage,
		adminMiddlewareService AdminMiddlewareService,
		handlers ...func(storage datastore.Storage, adminMiddlewareService AdminMiddlewareService) Handler,
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

func (s *server) AddFrontend(static string, index string) {
	s.Echo.Static("/", static)
	s.Echo.File("/", index)
}

func (s *server) AddHandlers(
	storage datastore.Storage,
	adminMiddlewareService AdminMiddlewareService,
	handlers ...func(storage datastore.Storage, adminMiddlewareService AdminMiddlewareService) Handler,
) {
	for _, handler := range handlers {
		s.addHandler(handler(storage, adminMiddlewareService))
	}
}

func (s *server) addHandler(h Handler) {
	h.AddRoutes(s.Echo.Group("/api"))
}

func (s *server) addGeneralMiddleware() {
	logging.AddMiddleware(s.Echo)

	s.Echo.Use(middleware.Recover())
	s.Echo.Use(middleware.Secure())

	port := viper.GetString("server.port")
	s.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{fmt.Sprintf("http://localhost:%s", port)},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))
}
