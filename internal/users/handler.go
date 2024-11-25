package users

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/server"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UsersHandler interface {
	server.Handler
	Register(c echo.Context) error
	Login(c echo.Context) error
	Refresh(c echo.Context) error
	Logout(c echo.Context) error
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
	publicGroup.POST("/register", h.Register)
	publicGroup.POST("/login", h.Login)
	publicGroup.POST("/refresh", h.Refresh)

	restrictedGroup := usersGroup.Group("")
	if err := h.authService.AddAuthMiddleware(restrictedGroup); err != nil {
		panic(err)
	}
	if err := h.usersService.AddUserPreloadMiddleware(restrictedGroup); err != nil {
		panic(err)
	}

	restrictedGroup.POST("/logout", h.Logout)
	restrictedGroup.GET("/me", h.GetMe)
}

func (h *UsersHandlerImpl) Register(c echo.Context) error {
	registerRequest := new(struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	})

	if err := c.Bind(registerRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(registerRequest); err != nil {
		return err
	}

	if err := h.authService.ValidatePassword(registerRequest.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.usersService.Create(registerRequest.Email)
	if err != nil && errors.Is(err, ErrUserAlreadyExists) {
		return echo.NewHTTPError(http.StatusConflict, "user already exists")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.authService.Register(user.ID, registerRequest.Password)
	if err != nil {
		slog.Error("failed to register user password, deleting user", slog.Uint64("user_id", uint64(user.ID)), slog.Any("error", err))
		if innerErr := h.usersService.Delete(user.ID); innerErr != nil {
			slog.Error("failed to delete user", slog.Uint64("user_id", uint64(user.ID)), slog.Any("error", innerErr))
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UsersHandlerImpl) Login(c echo.Context) error {
	loginRequest := new(struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	})

	if err := c.Bind(loginRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(loginRequest); err != nil {
		return err
	}

	user, err := h.usersService.GetByEmail(loginRequest.Email)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	tokenPair, err := h.authService.Login(user.ID, loginRequest.Password)
	if err != nil && errors.Is(err, auth.ErrInvalidPassword) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid email or password")
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
