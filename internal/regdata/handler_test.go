package regdata_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/L2SH-Dev/admissions/internal/regdata"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/L2SH-Dev/admissions/internal/validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	e *echo.Echo
	h regdata.RegistrationDataHandler
)

func setupTestHandler(t *testing.T) {
	t.Cleanup(func() {
		err := storage.Flush()
		assert.NoError(t, err)
	})

	e = echo.New()
	validation.AddValidation(e)
	h = regdata.NewRegistrationDataHandler(storage).(regdata.RegistrationDataHandler)

	// Create default roles for testing
	rolesRepo := roles.NewRolesRepo(storage)
	rolesService := roles.NewRolesService(rolesRepo)
	err := rolesService.CreateDefaultRoles()
	require.NoError(t, err)
}

func TestRegister(t *testing.T) {
	setupTestHandler(t)

	validData := regdata.RegistrationData{
		Email:           "test@example.com",
		FirstName:       "Test",
		LastName:        "User",
		Gender:          "M",
		BirthDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Grade:           9,
		OldSchool:       "Previous School",
		ParentFirstName: "Parent",
		ParentLastName:  "Test",
		ParentPhone:     "+1234567890",
	}

	// Test valid registration
	jsonData, err := json.Marshal(validData)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/regdata", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, h.Register(c))
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Test invalid data
	invalidData := regdata.RegistrationData{
		Email: "invalid", // Invalid email format
	}
	jsonData, err = json.Marshal(invalidData)
	require.NoError(t, err)

	req = httptest.NewRequest(http.MethodPost, "/regdata", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = h.Register(c)
	assert.Error(t, err)
}

func TestVerifyEmail(t *testing.T) {
	setupTestHandler(t)

	// First create a registration
	data := regdata.RegistrationData{
		Email:           "test@example.com",
		FirstName:       "Test",
		LastName:        "User",
		Gender:          "M",
		BirthDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Grade:           9,
		OldSchool:       "Previous School",
		ParentFirstName: "Parent",
		ParentLastName:  "Test",
		ParentPhone:     "+1234567890",
	}

	jsonData, err := json.Marshal(data)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/regdata", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Register(c))

	// Test invalid verification token
	req = httptest.NewRequest(http.MethodGet, "/regdata/verify/invalid-token", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("verification_token")
	c.SetParamValues("invalid-token")

	err = h.VerifyEmail(c)
	assert.Error(t, err)

	// Test empty verification token
	req = httptest.NewRequest(http.MethodGet, "/regdata/verify/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("verification_token")
	c.SetParamValues("")

	err = h.VerifyEmail(c)
	assert.Error(t, err)
}

func TestAccept(t *testing.T) {
	setupTestHandler(t)

	// First create and verify a registration
	data := regdata.RegistrationData{
		Email:           "test@example.com",
		FirstName:       "Test",
		LastName:        "User",
		Gender:          "M",
		BirthDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Grade:           9,
		OldSchool:       "Previous School",
		ParentFirstName: "Parent",
		ParentLastName:  "Test",
		ParentPhone:     "+1234567890",
	}

	jsonData, err := json.Marshal(data)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/regdata", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Register(c))

	var respData regdata.RegistrationData
	err = json.Unmarshal(rec.Body.Bytes(), &respData)
	require.NoError(t, err)

	// Test accepting unverified registration
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/regdata/admin/accept/%d", respData.ID), nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprint(respData.ID))

	err = h.Accept(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email is not verified")

	// Test accepting non-existent registration
	req = httptest.NewRequest(http.MethodPost, "/regdata/admin/accept/999", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")

	err = h.Accept(c)
	assert.Error(t, err)

	// Test invalid ID format
	req = httptest.NewRequest(http.MethodPost, "/regdata/admin/accept/invalid", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	err = h.Accept(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid registration data ID")
}

func TestListPending(t *testing.T) {
	setupTestHandler(t)

	// Create test data
	data := regdata.RegistrationData{
		Email:           "test@example.com",
		FirstName:       "Test",
		LastName:        "User",
		Gender:          "M",
		BirthDate:       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		Grade:           9,
		OldSchool:       "Previous School",
		ParentFirstName: "Parent",
		ParentLastName:  "Test",
		ParentPhone:     "+1234567890",
		EmailVerified:   true, // Verified from the start
	}

	jsonData, err := json.Marshal(data)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/regdata", bytes.NewBuffer(jsonData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	require.NoError(t, h.Register(c))

	// Test listing registrations
	req = httptest.NewRequest(http.MethodGet, "/regdata/admin", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	assert.NoError(t, h.ListPending(c))
	assert.Equal(t, http.StatusOK, rec.Code)

	var registrations []*regdata.RegistrationData
	err = json.Unmarshal(rec.Body.Bytes(), &registrations)
	require.NoError(t, err)
	assert.Len(t, registrations, 1)
	assert.Equal(t, data.Email, registrations[0].Email)
}
