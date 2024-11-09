package server

import (
	"fmt"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/database"
	"github.com/L2SH-Dev/admissions/internal/logging"
	"github.com/L2SH-Dev/admissions/internal/validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Server interface {
	Start()
	GetEcho() *echo.Echo
	GetDB() *gorm.DB
	GetCache() *redis.Client
}

type server struct {
	Echo  *echo.Echo
	DB    *gorm.DB
	Cache *redis.Client
}

func NewServer() Server {
	srv := &server{
		Echo:  echo.New(),
		DB:    database.InitDBConnection(),
		Cache: database.InitCacheConnection(),
	}

	srv.addMiddleware()
	validation.AddValidation(srv.Echo)

	return srv
}

func (s *server) Start() {
	port := viper.GetString("server.port")
	s.Echo.Logger.Fatal(s.Echo.Start(fmt.Sprintf(":%s", port)))
}

func (s *server) GetEcho() *echo.Echo {
	return s.Echo
}

func (s *server) GetDB() *gorm.DB {
	return s.DB
}

func (s *server) GetCache() *redis.Client {
	return s.Cache
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
