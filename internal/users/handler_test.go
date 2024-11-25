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

		err = storage.Cache.FlushDB(context.Background()).Err()
		assert.NoError(t, err)
	})

	return users.NewUsersHandler(storage).(users.UsersHandler)
}

func setupEcho() *echo.Echo {
	e := echo.New()
	validation.AddValidation(e)
	return e
}

func registerUser(t *testing.T, e *echo.Echo, handler users.UsersHandler) {
	req := httptest.NewRequest(http.MethodPost, "/users/register", strings.NewReader(`{"email":"test@example.com", "password":"Password$123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Register(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func loginUser(t *testing.T, e *echo.Echo, handler users.UsersHandler) (string, string) {
	req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"email":"test@example.com", "password":"Password$123"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var loginResponse struct {
		AccessToken  string `json:"access"`
		RefreshToken string `json:"refresh"`
	}
	err = json.NewDecoder(rec.Body).Decode(&loginResponse)
	assert.NoError(t, err)

	return loginResponse.AccessToken, loginResponse.RefreshToken
}

func TestUsersHandler_Register(t *testing.T) {
	handler := setupTestHandler(t)
	e := setupEcho()

	registerUser(t, e, handler)
}

func TestUsersHandler_Login(t *testing.T) {
	handler := setupTestHandler(t)
	e := setupEcho()

	registerUser(t, e, handler)
	loginUser(t, e, handler)
}

func TestUsersHandler_Refresh(t *testing.T) {
	handler := setupTestHandler(t)
	e := setupEcho()

	registerUser(t, e, handler)
	_, refreshToken := loginUser(t, e, handler)

	// Use refresh token to get new token pair
	req := httptest.NewRequest(http.MethodPost, "/users/refresh", strings.NewReader(`{"refresh":"`+refreshToken+`"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.Refresh(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"access"`)
	assert.Contains(t, rec.Body.String(), `"refresh"`)
}

func TestUsersHandler_GetMe(t *testing.T) {
	handler := setupTestHandler(t)
	e := setupEcho()

	handler.AddRoutes(e.Group(""))

	registerUser(t, e, handler)
	accessToken, _ := loginUser(t, e, handler)

	// Use access token to get user info
	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+accessToken)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"email":"test@example.com"`)
}

func TestUsersHandler_Logout(t *testing.T) {
	handler := setupTestHandler(t)
	e := setupEcho()

	handler.AddRoutes(e.Group(""))

	registerUser(t, e, handler)
	accessToken, _ := loginUser(t, e, handler)

	// Use access token to logout
	req := httptest.NewRequest(http.MethodPost, "/users/logout", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+accessToken)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"logged out"`)

	// Try to get user info again
	req = httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+accessToken)
	rec = httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), `"message":"token not found"`)
}
