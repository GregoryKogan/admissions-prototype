package auth

import (
	"log/slog"

	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func AddJWTMiddleware(g *echo.Group) error {
	jwtKey, err := secrets.ReadSecret("jwt_key")
	if err != nil {
		slog.Error("failed to add JWT middleware", slog.Any("error", err))
		return err
	}

	g.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTClaims)
		},
		SigningKey: []byte(jwtKey),
	}))

	return nil
}
