package ping

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PingHandler interface {
	AddRoutes(g *echo.Group)
	Ping(c echo.Context) error
}

type PingHandlerImpl struct{}

func NewPingHandler() PingHandler {
	return &PingHandlerImpl{}
}

func (h *PingHandlerImpl) AddRoutes(g *echo.Group) {
	g.GET("/ping", h.Ping)
}

func (h *PingHandlerImpl) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
