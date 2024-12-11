package exams

import (
	"net/http"
	"strconv"

	"github.com/L2SH-Dev/admissions/internal/datastore"
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
	ListAvailable(c echo.Context) error
	Register(c echo.Context) error
	Allocation(c echo.Context) error

	// admin endpoints
	List(c echo.Context) error
	Create(c echo.Context) error
	Delete(c echo.Context) error
	ListTypes(c echo.Context) error
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

	repo := NewExamsRepo(storage)
	service := NewExamsService(repo)

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

	privateGroup.GET("/available", h.ListAvailable)
	privateGroup.POST("/register/:examID", h.Register)
	privateGroup.GET("/allocation/:examID", h.Allocation)

	// admin endpoints
	adminGroup := privateGroup.Group("/admin")
	usersMiddlewareService.AddAdminMiddleware(adminGroup, roles.Role{WriteGeneral: true})

	adminGroup.GET("", h.List)
	adminGroup.POST("", h.Create)
	adminGroup.DELETE("/:examID", h.Delete)
	adminGroup.GET("/types", h.ListTypes)
}

func (h *ExamsHandlerImpl) ListAvailable(c echo.Context) error {
	user := c.Get("currentUser").(*users.User)
	exams, err := h.service.ListAvailable(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, exams)
}

func (h *ExamsHandlerImpl) Register(c echo.Context) error {
	user := c.Get("currentUser").(*users.User)
	examID64, err := strconv.ParseUint(c.Param("examID"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid exam ID")
	}
	examID := uint(examID64)

	err = h.service.Register(user, examID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (h *ExamsHandlerImpl) Allocation(c echo.Context) error {
	examID64, err := strconv.ParseUint(c.Param("examID"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid exam ID")
	}
	examID := uint(examID64)

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
	examID64, err := strconv.ParseUint(c.Param("examID"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid exam ID")
	}
	examID := uint(examID64)

	err = h.service.Delete(examID)
	if err != nil {
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
