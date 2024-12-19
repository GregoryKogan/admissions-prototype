package regdata

import (
	"bytes"
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"
	"time"
	_ "time/tzdata"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/mailing"
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

	// Private endpoints
	GetMine(c echo.Context) error

	// Admin endpoints
	Accept(c echo.Context) error
	Reject(c echo.Context) error
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

	// Private endpoints
	privateGroup := regDataGroup.Group("")

	usersMiddlewareService := users.NewUsersMiddlewareService(h.usersService, h.authService)
	jwtKey := viper.GetString("secrets.jwt_key")
	usersMiddlewareService.AddAuthMiddleware(privateGroup, jwtKey)
	usersMiddlewareService.AddUserPreloadMiddleware(privateGroup)

	privateGroup.GET("/mine", h.GetMine)

	// Admin endpoints
	adminGroup := privateGroup.Group("/admin")
	usersMiddlewareService.AddAdminMiddleware(adminGroup, roles.Role{WriteGeneral: true})

	adminGroup.POST("/accept/:id", h.Accept)
	adminGroup.POST("/reject/:id", h.Reject)
	adminGroup.GET("/pending", h.ListPending)
	adminGroup.GET("/accepted", h.ListAccepted)
	adminGroup.GET("/accepted/download", h.DownloadAcceptedRegistrations)
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

func (h *RegistrationDataHandlerImpl) GetMine(c echo.Context) error {
	user := c.Get("currentUser").(*users.User)
	regData, err := h.service.GetByID(user.RegistrationDataID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, regData)
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

func (h *RegistrationDataHandlerImpl) Reject(c echo.Context) error {
	regDataID64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid registration data ID")
	}
	regDataID := uint(regDataID64)

	rejectRequest := new(struct {
		Reason string `json:"reason" validate:"required"`
	})

	if err := c.Bind(rejectRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(rejectRequest); err != nil {
		return err
	}

	regData, err := h.service.GetByID(regDataID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	email := regData.Email

	err = h.service.Reject(regDataID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = mailing.SendRegistrationRejection(email, rejectRequest.Reason)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *RegistrationDataHandlerImpl) ListPending(c echo.Context) error {
	registrations, err := h.service.GetPending()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, registrations)
}

func (h *RegistrationDataHandlerImpl) ListAccepted(c echo.Context) error {
	registrations, err := h.service.GetAccepted()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, registrations)
}

func (h *RegistrationDataHandlerImpl) DownloadAcceptedRegistrations(c echo.Context) error {
	registrations, err := h.service.GetAccepted()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Prepare CSV data with UTF-8 BOM
	var csvData bytes.Buffer
	// Write UTF-8 BOM
	csvData.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(&csvData)
	// Write header
	writer.Write([]string{
		"№",
		"ID",
		"Email",
		"Фамилия",
		"Имя",
		"Отчество",
		"Дата рождения",
		"Класс поступления",
		"Телефон родителя",
		"Фамилия родителя",
		"Имя родителя",
		"Отчество родителя",
		"Школа",
		"ВМШ",
		"Июньский экзамен",
		"Как узнали о Лицее",
		"Логин",
		"Дата регистрации",
	})

	tz, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Write data with line numbers
	for i, reg := range registrations {
		writer.Write([]string{
			strconv.Itoa(i + 1),
			strconv.FormatUint(uint64(reg.ID), 10),
			reg.Email,
			reg.LastName,
			reg.FirstName,
			reg.Patronymic,
			reg.BirthDate.In(tz).Format("2006.01.02"),
			strconv.FormatUint(uint64(reg.Grade), 10),
			reg.ParentPhone,
			reg.ParentLastName,
			reg.ParentFirstName,
			reg.ParentPatronymic,
			reg.OldSchool,
			strconv.FormatBool(reg.VMSH),
			strconv.FormatBool(reg.JuneExam),
			reg.Source,
			reg.User.Login,
			reg.CreatedAt.In(tz).Format("2006.01.02 15:04:05"),
		})
	}
	writer.Flush()

	return c.Blob(http.StatusOK, "text/csv; charset=utf-8", csvData.Bytes())
}
