package users

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/auth"
	"github.com/L2SH-Dev/admissions/internal/passwords"
	"github.com/L2SH-Dev/admissions/internal/server"
	"github.com/L2SH-Dev/admissions/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

type UsersHandler interface {
	server.Handler
	Register(c echo.Context) error
	Login(c echo.Context) error
	Refresh(c echo.Context) error
	GetMe(c echo.Context) error
}

type UsersHandlerImpl struct {
	usersService UsersService
	authService  auth.AuthService
}

func NewUsersHandler(storage storage.Storage) server.Handler {
	usersRepo := NewUsersRepo(storage)
	usersService := NewUsersService(usersRepo)
	passwordsRepo := passwords.NewPasswordsRepo(storage)
	passwordsService := passwords.NewPasswordsService(passwordsRepo)
	authService := auth.NewAuthService(passwordsService)
	return &UsersHandlerImpl{
		usersService: usersService,
		authService:  authService,
	}
}

func (h *UsersHandlerImpl) AddRoutes(g *echo.Group) {
	usersGroup := g.Group("/users")

	authGroup := usersGroup.Group("/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)
	authGroup.POST("/refresh", h.Refresh)

	restrictedGroup := usersGroup.Group("")
	if err := auth.AddJWTMiddleware(restrictedGroup); err != nil {
		panic(err)
	}

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
	fmt.Println(c.Get("user"))
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*auth.JWTClaims)
	userID := claims.UserID

	user, err := h.usersService.GetFullByID(userID)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}
