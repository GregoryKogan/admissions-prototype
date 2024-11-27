package users

import (
	"errors"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UsersMiddlewareService interface {
	AddAuthMiddleware(g *echo.Group, jwtKey string)
	AddUserPreloadMiddleware(g *echo.Group)
	AddAdminMiddleware(g *echo.Group, minimalRole roles.Role)
}

type UsersMiddlewareServiceImpl struct {
	usersService UsersService
	authService  auth.AuthService
}

func NewUsersMiddlewareService(usersService UsersService, authService auth.AuthService) UsersMiddlewareService {
	return &UsersMiddlewareServiceImpl{
		usersService: usersService,
		authService:  authService,
	}
}

func (s *UsersMiddlewareServiceImpl) AddAuthMiddleware(g *echo.Group, jwtKey string) {
	s.authService.AddAuthMiddleware(g, jwtKey)
}

func (s *UsersMiddlewareServiceImpl) AddUserPreloadMiddleware(g *echo.Group) {
	g.Use(s.preloadUserDataMiddleware())
}

func (s *UsersMiddlewareServiceImpl) AddAdminMiddleware(g *echo.Group, minimalRole roles.Role) {
	g.Use(userAdminMiddleware(minimalRole))
}

func (s *UsersMiddlewareServiceImpl) preloadUserDataMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("userId").(uint)

			userDetails, err := s.usersService.GetByID(userID)
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

func userAdminMiddleware(minimalRole roles.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("currentUser").(*User)
			if !user.Role.Admin {
				return echo.NewHTTPError(http.StatusForbidden, "only admins can access this endpoint")
			}

			if !user.Role.WriteGeneral && minimalRole.WriteGeneral {
				return echo.NewHTTPError(http.StatusForbidden, "only admins with general write permissions can access this endpoint")
			}

			if !user.Role.AIAccess && minimalRole.AIAccess {
				return echo.NewHTTPError(http.StatusForbidden, "only admins with AI access can access this endpoint")
			}

			return next(c)
		}
	}
}
