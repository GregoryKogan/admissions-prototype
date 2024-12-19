package exams

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/regdata"
	"github.com/L2SH-Dev/admissions/internal/server"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type ExamsHandler interface {
	server.Handler

	// private endpoints
	History(c echo.Context) error
	Available(c echo.Context) error
	Register(c echo.Context) error
	Unregister(c echo.Context) error
	Allocation(c echo.Context) error
	RegistrationStatus(c echo.Context) error

	// admin endpoints
	List(c echo.Context) error
	Create(c echo.Context) error
	Delete(c echo.Context) error
	ListTypes(c echo.Context) error
	DownloadRegistrations(c echo.Context) error
}

type ExamsHandlerImpl struct {
	service      ExamsService
	usersService users.UsersService
	authService  auth.AuthService
}

func NewExamsHandler(storage datastore.Storage) server.Handler {
	rolesRepo := roles.NewRolesRepo(storage)
	rolesService := roles.NewRolesService(rolesRepo)
	usersRepo := users.NewUsersRepo(storage)
	usersService := users.NewUsersService(usersRepo, rolesService)

	passwordsRepo := passwords.NewPasswordsRepo(storage)
	passwordsService := passwords.NewPasswordsService(passwordsRepo)
	authRepo := auth.NewAuthRepo(storage)
	authService := auth.NewAuthService(authRepo, passwordsService)

	regDataRepo := regdata.NewRegistrationDataRepo(storage)
	regDataService := regdata.NewRegistrationDataService(regDataRepo, usersService, authService, passwordsService)

	repo := NewExamsRepo(storage)
	service := NewExamsService(repo, regDataService)

	service.CreateDefaultExamTypes()

	return &ExamsHandlerImpl{
		service:      service,
		usersService: usersService,
		authService:  authService,
	}
}

func (h *ExamsHandlerImpl) AddRoutes(g *echo.Group) {
	examsGroup := g.Group("/exams")

	// private endpoints
	privateGroup := examsGroup.Group("")
	usersMiddlewareService := users.NewUsersMiddlewareService(h.usersService, h.authService)
	jwtKey := viper.GetString("secrets.jwt_key")
	usersMiddlewareService.AddAuthMiddleware(privateGroup, jwtKey)
	usersMiddlewareService.AddUserPreloadMiddleware(privateGroup)

	privateGroup.GET("/history", h.History)
	privateGroup.GET("/available", h.Available)
	privateGroup.POST("/register/:examID", h.Register)
	privateGroup.DELETE("/register/:examID", h.Unregister)
	privateGroup.GET("/allocation/:examID", h.Allocation)
	privateGroup.GET("/registration_status/:examID", h.RegistrationStatus)

	// admin endpoints
	adminGroup := privateGroup.Group("/admin")
	usersMiddlewareService.AddAdminMiddleware(adminGroup, roles.Role{WriteGeneral: true})

	adminGroup.GET("", h.List)
	adminGroup.POST("", h.Create)
	adminGroup.DELETE("/:examID", h.Delete)
	adminGroup.GET("/types", h.ListTypes)
	adminGroup.GET("/registrations/:examID/download", h.DownloadRegistrations)
}

func (h *ExamsHandlerImpl) History(c echo.Context) error {
	user := c.Get("currentUser").(*users.User)
	exams, err := h.service.History(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, exams)
}

func (h *ExamsHandlerImpl) Available(c echo.Context) error {
	user := c.Get("currentUser").(*users.User)
	exams, err := h.service.Available(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, exams)
}

func (h *ExamsHandlerImpl) Register(c echo.Context) error {
	user := c.Get("currentUser").(*users.User)
	examID, err := parseUintParam(c, "examID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid exam ID")
	}

	if err := h.service.Register(user, examID); err != nil {
		return mapServiceError(err)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *ExamsHandlerImpl) Unregister(c echo.Context) error {
	user := c.Get("currentUser").(*users.User)
	examID, err := parseUintParam(c, "examID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid exam ID")
	}

	if err := h.service.Unregister(user, examID); err != nil {
		return mapServiceError(err)
	}

	return c.NoContent(http.StatusOK)
}

func (h *ExamsHandlerImpl) Allocation(c echo.Context) error {
	examID, err := parseUintParam(c, "examID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid exam ID")
	}

	allocation, err := h.service.Allocation(examID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, allocation)
}

func (h *ExamsHandlerImpl) List(c echo.Context) error {
	exams, err := h.service.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, exams)
}

func (h *ExamsHandlerImpl) Create(c echo.Context) error {
	exam := new(Exam)
	if err := c.Bind(exam); err != nil {
		return err
	}

	if err := c.Validate(exam); err != nil {
		return err
	}

	err := h.service.Create(exam)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, exam)
}

func (h *ExamsHandlerImpl) Delete(c echo.Context) error {
	examID, err := parseUintParam(c, "examID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid exam ID")
	}

	if err := h.service.Delete(examID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *ExamsHandlerImpl) ListTypes(c echo.Context) error {
	types, err := h.service.ListTypes()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, types)
}

func (h *ExamsHandlerImpl) RegistrationStatus(c echo.Context) error {
	user := c.Get("currentUser").(*users.User)
	examID, err := parseUintParam(c, "examID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid exam ID")
	}

	registeredToExam, registeredToSameType, err := h.service.RegistrationStatus(user, examID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]bool{
		"registered":              registeredToExam,
		"registered_to_same_type": registeredToSameType,
	})
}

func (h *ExamsHandlerImpl) DownloadRegistrations(c echo.Context) error {
	examID, err := parseUintParam(c, "examID")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid exam ID")
	}

	registrations, err := h.service.GetRegistrations(examID)
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
		"Фамилия",
		"Имя",
		"Отчество",
		"Телефон",
		"Фамилия родителя",
		"Имя родителя",
		"Отчество родителя",
	})
	// Write data with line numbers
	for i, reg := range registrations {
		writer.Write([]string{
			strconv.Itoa(i + 1),
			strconv.FormatUint(uint64(reg.User.ID), 10),
			reg.LastName,
			reg.FirstName,
			reg.Patronymic,
			reg.ParentPhone,
			reg.ParentLastName,
			reg.ParentFirstName,
			reg.ParentPatronymic,
		})
	}
	writer.Flush()

	return c.Blob(http.StatusOK, "text/csv; charset=utf-8", csvData.Bytes())
}

func parseUintParam(c echo.Context, param string) (uint, error) {
	value64, err := strconv.ParseUint(c.Param(param), 10, 32)
	return uint(value64), err
}

func mapServiceError(err error) *echo.HTTPError {
	switch err {
	case ErrAlreadyRegistered, ErrInvalidGrade, ErrExamFull, ErrInvalidExamOrder:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	case ErrNotRegistered:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
