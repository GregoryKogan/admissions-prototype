package users

import (
	"errors"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/auth"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

type UsersHandler struct {
	usersService *UsersService
	authService  *auth.AuthService
}

func NewUsersHandler(usersService *UsersService, authService *auth.AuthService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
		authService:  authService,
	}
}

func (h *UsersHandler) AddRoutes(g *echo.Group) {
	usersGroup := g.Group("/users")
	authGroup := usersGroup.Group("/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)
}

func (h *UsersHandler) Register(c echo.Context) error {
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

func (h *UsersHandler) Login(c echo.Context) error {
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
