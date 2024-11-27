package users

import (
	"errors"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/server"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UsersHandler interface {
	server.Handler
	Login(c echo.Context) error
	Logout(c echo.Context) error
	Refresh(c echo.Context) error
	GetMe(c echo.Context) error
}

type UsersHandlerImpl struct {
	usersService UsersService
	authService  auth.AuthService
}

func NewUsersHandler(storage datastore.Storage) server.Handler {
	rolesRepo := roles.NewRolesRepo(storage)
	rolesService := roles.NewRolesService(rolesRepo)
	usersRepo := NewUsersRepo(storage)
	usersService := NewUsersService(usersRepo, rolesService)

	passwordsRepo := passwords.NewPasswordsRepo(storage)
	passwordsService := passwords.NewPasswordsService(passwordsRepo)

	authRepo := auth.NewAuthRepo(storage)
	authService := auth.NewAuthService(authRepo, passwordsService)

	return &UsersHandlerImpl{
		usersService: usersService,
		authService:  authService,
	}
}

func (h *UsersHandlerImpl) AddRoutes(g *echo.Group) {
	usersGroup := g.Group("/users")

	publicGroup := usersGroup.Group("")
	publicGroup.POST("/login", h.Login)
	publicGroup.POST("/refresh", h.Refresh)

	restrictedGroup := usersGroup.Group("")

	middlewareService := NewUsersMiddlewareService(h.usersService, h.authService)
	middlewareService.AddAuthMiddleware(restrictedGroup, viper.GetString("secrets.jwt_key"))
	middlewareService.AddUserPreloadMiddleware(restrictedGroup)

	restrictedGroup.POST("/logout", h.Logout)
	restrictedGroup.GET("/me", h.GetMe)
}

func (h *UsersHandlerImpl) Login(c echo.Context) error {
	loginRequest := new(struct {
		Login    string `json:"login" validate:"required"`
		Password string `json:"password" validate:"required"`
	})

	if err := c.Bind(loginRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(loginRequest); err != nil {
		return err
	}

	user, err := h.usersService.GetByLogin(loginRequest.Login)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid login or password")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tokenPair, err := h.authService.Login(user.ID, loginRequest.Password)
	if err != nil && errors.Is(err, auth.ErrInvalidPassword) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid login or password")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tokenPair)
}

func (h *UsersHandlerImpl) Logout(c echo.Context) error {
	user := c.Get("currentUser").(*User)
	h.authService.Logout(user.ID)
	return c.JSON(http.StatusOK, "logged out")
}

func (h *UsersHandlerImpl) Refresh(c echo.Context) error {
	refreshRequest := new(struct {
		RefreshToken string `json:"refresh" validate:"required"`
	})

	if err := c.Bind(refreshRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(refreshRequest); err != nil {
		return err
	}

	tokenPair, err := h.authService.Refresh(refreshRequest.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, tokenPair)
}

func (h *UsersHandlerImpl) GetMe(c echo.Context) error {
	user := c.Get("currentUser").(*User)
	return c.JSON(http.StatusOK, user)
}
