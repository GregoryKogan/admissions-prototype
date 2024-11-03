package ping

import "github.com/labstack/echo/v4"

func AddRoutes(g *echo.Group) {
	g.GET("/ping", PingHandler)
}
