package auth

import (
	"github.com/labstack/echo/v4"
)

func AddRoutes(g *echo.Group) {
	authGroup := g.Group("/auth")
	authGroup.POST("/login", LoginHandler)
	authGroup.POST("/register", RegisterHandler)
	authGroup.POST("/refresh", RefreshHandler)
	authGroup.POST("/logout", LogoutHandler)
}
