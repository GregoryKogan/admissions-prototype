package validation

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Validator interface {
	Validate(i interface{}) error
}

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() Validator {
	return &CustomValidator{validator: validator.New()}
}

func AddValidation(e *echo.Echo) {
	e.Validator = NewCustomValidator()
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
