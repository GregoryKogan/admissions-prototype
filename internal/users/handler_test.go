package users_test

// import (
// 	"testing"

// 	"github.com/L2SH-Dev/admissions/internal/server"
// 	"github.com/L2SH-Dev/admissions/internal/users"
// 	"github.com/L2SH-Dev/admissions/internal/validation"
// 	"github.com/labstack/echo/v4"
// 	"github.com/stretchr/testify/assert"
// )

// var (
// 	e       *echo.Echo
// 	handler server.Handler
// )

// func setupTestHandler(t *testing.T) {
// 	t.Cleanup(func() {
// 		err := storage.Flush()
// 		assert.NoError(t, err)
// 	})

// 	handler = users.NewUsersHandler(storage)
// 	e = echo.New()
// 	validation.AddValidation(e)
// 	handler.AddRoutes(e.Group(""))
// }
