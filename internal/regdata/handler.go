package regdata

import (
	"errors"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"github.com/L2SH-Dev/admissions/internal/regdata/emailver"
	"github.com/L2SH-Dev/admissions/internal/server"
	"github.com/labstack/echo/v4"
)

type RegistrationDataHandler interface {
	server.Handler
	Register(c echo.Context) error
	VerifyEmail(c echo.Context) error
}

type RegistrationDataHandlerImpl struct {
	service                  RegistrationDataService
	emailVerificationService emailver.EmailVerificationService
}

func NewRegistrationDataHandler(storage datastore.Storage) server.Handler {
	repo := NewRegistrationDataRepo(storage)
	service := NewRegistrationDataService(repo)

	emailVerRepo := emailver.NewEmailVerificationRepo(storage)
	emailVerService := emailver.NewEmailVerificationService(emailVerRepo)

	return &RegistrationDataHandlerImpl{service: service, emailVerificationService: emailVerService}
}

func (h *RegistrationDataHandlerImpl) AddRoutes(g *echo.Group) {
	regDataGroup := g.Group("/regdata")
	regDataGroup.POST("", h.Register)
	regDataGroup.POST("verify/:verification_token", h.VerifyEmail)
}

func (h *RegistrationDataHandlerImpl) Register(c echo.Context) error {
	data := new(RegistrationData)
	if err := c.Bind(data); err != nil {
		return err
	}

	if err := c.Validate(data); err != nil {
		return err
	}

	err := h.service.CreateRegistrationData(data)
	if err != nil && errors.Is(err, ErrRegistrationDataInvalid) {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	} else if err != nil && errors.Is(err, ErrorRegistrationDataExists) {
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
	verificationToken := c.Param("verification_token")
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

	return c.JSON(http.StatusOK, nil)
}
