package users

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (s *UsersServiceImpl) AddUserPreloadMiddleware(g *echo.Group) error {
	g.Use(s.preloadUserDataMiddleware())
	return nil
}

func (s *UsersServiceImpl) preloadUserDataMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("userId").(uint)

			userDetails, err := s.GetByID(userID)
			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				return echo.NewHTTPError(http.StatusNotFound, "user not found")
			} else if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			c.Set("currentUser", userDetails)

			return next(c)
		}
	}
}
