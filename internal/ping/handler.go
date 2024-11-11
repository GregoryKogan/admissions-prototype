package ping

import (
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/server"
	"github.com/labstack/echo/v4"
)

type PingHandler interface {
	server.Handler
	Ping(c echo.Context) error
}

type PingHandlerImpl struct{}

func NewPingHandler(_ datastore.Storage) server.Handler {
	return &PingHandlerImpl{}
}

func (h *PingHandlerImpl) AddRoutes(g *echo.Group) {
	g.GET("/ping", h.Ping)
}

func (h *PingHandlerImpl) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
