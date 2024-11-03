package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	return c.String(http.StatusOK, "login")
}

func RegisterHandler(c echo.Context) error {
	return c.String(http.StatusOK, "register")
}

func RefreshHandler(c echo.Context) error {
	return c.String(http.StatusOK, "refresh")
}

func LogoutHandler(c echo.Context) error {
	return c.String(http.StatusOK, "logout")
}
