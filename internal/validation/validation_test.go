package validation_test

import (
	"net/http"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `validate:"required"`
}

func TestAddValidation(t *testing.T) {
	e := echo.New()
	validation.AddValidation(e)

	assert.NotNil(t, e.Validator)
}

func TestCustomValidator_Validate(t *testing.T) {
	e := echo.New()
	validation.AddValidation(e)

	validStruct := TestStruct{Name: "Valid Name"}
	invalidStruct := TestStruct{Name: ""}

	// Test valid struct
	err := e.Validator.Validate(validStruct)
	assert.NoError(t, err)

	// Test invalid struct
	err = e.Validator.Validate(invalidStruct)
	httpError, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpError.Code)
	assert.Contains(t, httpError.Message, "Name")
}
