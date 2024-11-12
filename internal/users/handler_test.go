package users_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestHandler(t *testing.T) users.UsersHandler {
	t.Cleanup(func() {
		err := storage.DB.Exec("DELETE FROM passwords").Error
		assert.NoError(t, err)

		err = storage.DB.Exec("DELETE FROM users").Error
		assert.NoError(t, err)

		err = storage.DB.Exec("DELETE FROM roles").Error
		assert.NoError(t, err)

		err = storage.Cache.FlushDB(context.Background()).Err()
		assert.NoError(t, err)
	})

	return users.NewUsersHandler(storage).(users.UsersHandler)
}

func TestUsersHandler_Register(t *testing.T) {
	handler := setupTestHandler(t)
	e := echo.New()
	validation.AddValidation(e)

	req := httptest.NewRequest(http.MethodPost, "/users/auth/register", strings.NewReader(`{"email":"test@example.com", "password":"Password$123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Register(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"email":"test@example.com"`)
}

func TestUsersHandler_Login(t *testing.T) {
	handler := setupTestHandler(t)
	e := echo.New()
	validation.AddValidation(e)

	// Register user first via endpoint
	req := httptest.NewRequest(http.MethodPost, "/users/auth/register", strings.NewReader(`{"email":"test@example.com", "password":"Password$123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Register(c)
	assert.NoError(t, err)

	// Login user
	req = httptest.NewRequest(http.MethodPost, "/users/auth/login", strings.NewReader(`{"email":"test@example.com", "password":"Password$123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = handler.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"access"`)
	assert.Contains(t, rec.Body.String(), `"refresh"`)
}

func TestUsersHandler_Refresh(t *testing.T) {
	handler := setupTestHandler(t)
	e := echo.New()
	validation.AddValidation(e)

	// Register user first via endpoint
	req := httptest.NewRequest(http.MethodPost, "/users/auth/register", strings.NewReader(`{"email":"test@example.com", "password":"Password$123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Register(c)
	assert.NoError(t, err)

	// Login user to get refresh token
	req = httptest.NewRequest(http.MethodPost, "/users/auth/login", strings.NewReader(`{"email":"test@example.com", "password":"Password$123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = handler.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Extract refresh token from response
	var loginResponse struct {
		RefreshToken string `json:"refresh"`
	}
	err = json.NewDecoder(rec.Body).Decode(&loginResponse)
	assert.NoError(t, err)

	// Use refresh token to get new token pair
	req = httptest.NewRequest(http.MethodPost, "/users/auth/refresh", strings.NewReader(`{"refresh":"`+loginResponse.RefreshToken+`"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	err = handler.Refresh(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"access"`)
	assert.Contains(t, rec.Body.String(), `"refresh"`)
}

func TestUsersHandler_GetMe(t *testing.T) {
	handler := setupTestHandler(t)
	e := echo.New()
	validation.AddValidation(e)

	handler.AddRoutes(e.Group(""))

	// Register user first via endpoint
	req := httptest.NewRequest(http.MethodPost, "/users/auth/register", strings.NewReader(`{"email":"test@example.com", "password":"Password$123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Login user to get access token
	req = httptest.NewRequest(http.MethodPost, "/users/auth/login", strings.NewReader(`{"email":"test@example.com", "password":"Password$123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Extract access token from response
	var loginResponse struct {
		AccessToken string `json:"access"`
	}
	err := json.NewDecoder(rec.Body).Decode(&loginResponse)
	assert.NoError(t, err)

	// Use access token to get user info
	req = httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+loginResponse.AccessToken)
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"email":"test@example.com"`)
}
