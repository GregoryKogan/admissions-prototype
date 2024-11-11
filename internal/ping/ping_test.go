package ping_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/ping"
	"github.com/L2SH-Dev/admissions/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestHandler(t *testing.T) ping.PingHandler {
	storage := storage.SetupMockStorage(t)
	return ping.NewPingHandler(storage).(ping.PingHandler)
}

func TestPingHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := setupTestHandler(t)
	if assert.NoError(t, handler.Ping(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "pong", rec.Body.String())
	}
}
