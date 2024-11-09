package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/auth"
	"github.com/L2SH-Dev/admissions/internal/passwords"
	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestHandler(t *testing.T) (users.UsersHandler, users.UsersService) {
	db := setupTestDB(t) // Reuse setupTestDB from repo_test.go
	repo := users.NewUsersRepo(db)
	service := users.NewUsersService(repo)
	passwordsRepo := passwords.NewPasswordsRepo(db)
	passwordsService := passwords.NewPasswordsService(passwordsRepo)
	authService := auth.NewAuthService(passwordsService)
	handler := users.NewUsersHandler(service, authService)
	return handler, service
}

func TestUsersHandler_Register(t *testing.T) {
	handler, _ := setupTestHandler(t)
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
	handler, _ := setupTestHandler(t)
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

	// Mock JWT key
	secrets.SetMockSecret("jwt_key", "test_key")
	defer secrets.ClearMockSecrets()

	err = handler.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"access"`)
	assert.Contains(t, rec.Body.String(), `"refresh"`)
}

func TestUsersHandler_Refresh(t *testing.T) {
	handler, _ := setupTestHandler(t)
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

	// Mock JWT key
	secrets.SetMockSecret("jwt_key", "test_key")
	defer secrets.ClearMockSecrets()

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
