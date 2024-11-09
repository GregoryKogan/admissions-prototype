package ping_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/ping"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := ping.NewPingHandler()
	if assert.NoError(t, handler.Ping(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "pong", rec.Body.String())
	}
}
