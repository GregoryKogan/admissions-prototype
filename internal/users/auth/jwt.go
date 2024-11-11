package auth

import (
	"errors"
	"time"

	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

var ErrInvalidToken = errors.New("invalid token")

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func NewAccessToken(userID uint) (string, error) {
	claims := newAccessJWTClaims(userID)
	return newSignedJWT(claims)
}

func NewRefreshToken(userID uint) (string, error) {
	claims := newRefreshJWTClaims(userID)
	return newSignedJWT(claims)
}

func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		jwtKey, err := secrets.ReadSecret("jwt_key")
		if err != nil {
			return nil, errors.Join(errors.New("failed to read jwt_key secret"), err)
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, errors.Join(ErrInvalidToken, err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func newSignedJWT(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey, err := secrets.ReadSecret("jwt_key")
	if err != nil {
		return "", errors.Join(errors.New("failed to read jwt_key secret"), err)
	}
	return token.SignedString([]byte(jwtKey))
}

func newAccessJWTClaims(userID uint) *JWTClaims {
	return &JWTClaims{
		UserID: userID,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.access_lifetime"))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "l2sh.admissions",
		},
	}
}

func newRefreshJWTClaims(userID uint) *JWTClaims {
	return &JWTClaims{
		UserID: userID,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.refresh_lifetime"))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "l2sh.admissions",
		},
	}
}
