package ping

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) AddRoutes(g *echo.Group) {
	g.GET("/ping", h.Ping)
}

func (h *PingHandler) Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
