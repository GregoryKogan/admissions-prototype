package auth

import (
	"log/slog"
	"net/http"

	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/L2SH-Dev/admissions/internal/users/auth/authjwt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func (s *AuthServiceImpl) AddAuthMiddleware(g *echo.Group) error {
	jwtKey, err := secrets.ReadSecret("jwt_key")
	if err != nil {
		slog.Error("failed to add JWT middleware", slog.Any("error", err))
		return err
	}

	g.Use(
		echojwt.WithConfig(
			echojwt.Config{
				NewClaimsFunc: func(c echo.Context) jwt.Claims {
					return new(authjwt.JWTClaims)
				},
				SigningKey: []byte(jwtKey),
			}),

		s.validateJWTMiddleware(),
	)

	return nil
}

func (s *AuthServiceImpl) validateJWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*authjwt.JWTClaims)

			if claims.Type != "access" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token type")
			}

			ok, err := s.IsTokenCached(claims)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "failed to validate token")
			}

			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "token not found")
			}

			c.Set("userId", claims.UserID)

			return next(c)
		}
	}
}
