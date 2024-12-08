package regdata

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/regdata/emailver"
	"github.com/L2SH-Dev/admissions/internal/server"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type RegistrationDataHandler interface {
	server.Handler
	Register(c echo.Context) error
	VerifyEmail(c echo.Context) error

	// Admin endpoints
	Accept(c echo.Context) error
	ListPending(c echo.Context) error
}

type RegistrationDataHandlerImpl struct {
	service                  RegistrationDataService
	emailVerificationService emailver.EmailVerificationService
	usersService             users.UsersService
	authService              auth.AuthService
}

func NewRegistrationDataHandler(storage datastore.Storage) server.Handler {
	emailVerRepo := emailver.NewEmailVerificationRepo(storage)
	emailVerService := emailver.NewEmailVerificationService(emailVerRepo)

	rolesRepo := roles.NewRolesRepo(storage)
	rolesService := roles.NewRolesService(rolesRepo)
	usersRepo := users.NewUsersRepo(storage)
	usersService := users.NewUsersService(usersRepo, rolesService)

	passwordsRepo := passwords.NewPasswordsRepo(storage)
	passwordsService := passwords.NewPasswordsService(passwordsRepo)
	authRepo := auth.NewAuthRepo(storage)
	authService := auth.NewAuthService(authRepo, passwordsService)

	repo := NewRegistrationDataRepo(storage)
	service := NewRegistrationDataService(repo, usersService, authService, passwordsService)

	return &RegistrationDataHandlerImpl{
		service:                  service,
		emailVerificationService: emailVerService,
		usersService:             usersService,
		authService:              authService,
	}
}

func (h *RegistrationDataHandlerImpl) AddRoutes(g *echo.Group) {
	regDataGroup := g.Group("/regdata")
	regDataGroup.POST("", h.Register)
	regDataGroup.GET("/verify", h.VerifyEmail)

	// Admin endpoints
	adminGroup := regDataGroup.Group("/admin")

	usersMiddlewareService := users.NewUsersMiddlewareService(h.usersService, h.authService)

	jwtKey := viper.GetString("secrets.jwt_key")
	usersMiddlewareService.AddAuthMiddleware(adminGroup, jwtKey)
	usersMiddlewareService.AddUserPreloadMiddleware(adminGroup)
	usersMiddlewareService.AddAdminMiddleware(adminGroup, roles.Role{WriteGeneral: true})

	adminGroup.POST("/accept/:id", h.Accept)
	adminGroup.GET("/pending", h.ListPending)
}

func (h *RegistrationDataHandlerImpl) Register(c echo.Context) error {
	data := new(RegistrationData)
	if err := c.Bind(data); err != nil {
		return err
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	err := h.service.Create(data)
	if err != nil && errors.Is(err, ErrRegistrationDataInvalid) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil && errors.Is(err, ErrRegistrationDataExists) {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.emailVerificationService.SendVerificationEmail(data.Email, data.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, data)
}

func (h *RegistrationDataHandlerImpl) VerifyEmail(c echo.Context) error {
	verificationToken := c.QueryParam("token")
	if verificationToken == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "verification token is required")
	}

	registrationID, err := h.emailVerificationService.VerifyEmail(verificationToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.service.SetEmailVerified(registrationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]bool{"verified": true})
}

func (h *RegistrationDataHandlerImpl) Accept(c echo.Context) error {
	regDataID64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid registration data ID")
	}
	regDataID := uint(regDataID64)

	user, err := h.service.Accept(regDataID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *RegistrationDataHandlerImpl) ListPending(c echo.Context) error {
	registrations, err := h.service.GetPending()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, registrations)
}
