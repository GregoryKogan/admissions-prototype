package admin

import (
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/server"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/labstack/echo/v4"
)

type AdminMiddlewareServiceImpl struct {
	storage datastore.Storage
}

func NewAdminMiddlewareService(storage datastore.Storage) server.AdminMiddlewareService {
	return &AdminMiddlewareServiceImpl{storage: storage}
}

func (s *AdminMiddlewareServiceImpl) Add(g *echo.Group, minimalRole roles.Role) error {
	usersService := initUsersService(s.storage)
	authService := initAuthService(s.storage)

	if err := authService.AddAuthMiddleware(g); err != nil {
		return err
	}

	if err := usersService.AddUserPreloadMiddleware(g); err != nil {
		return err
	}

	g.Use(userAdminMiddleware(minimalRole))

	return nil
}

func userAdminMiddleware(minimalRole roles.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("currentUser").(*users.User)
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

func initUsersService(storage datastore.Storage) users.UsersService {
	rolesRepo := roles.NewRolesRepo(storage)
	rolesService := roles.NewRolesService(rolesRepo)
	usersRepo := users.NewUsersRepo(storage)
	return users.NewUsersService(usersRepo, rolesService)
}

func initAuthService(storage datastore.Storage) auth.AuthService {
	passwordsRepo := passwords.NewPasswordsRepo(storage)
	passwordsService := passwords.NewPasswordsService(passwordsRepo)
	authRepo := auth.NewAuthRepo(storage)
	return auth.NewAuthService(authRepo, passwordsService)
}
