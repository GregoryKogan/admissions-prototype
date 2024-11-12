package authjwt

import (
	"errors"
	"time"

	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var ErrInvalidToken = errors.New("invalid token")

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	UID    string `json:"uid"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

type JWTService interface {
	NewAccessToken(userID uint) (string, error)
	NewRefreshToken(userID uint) (string, error)
	ParseToken(tokenString string) (*JWTClaims, error)
}

type JWTServiceImpl struct{}

func NewJWTService() JWTService {
	return &JWTServiceImpl{}
}

func (s *JWTServiceImpl) NewAccessToken(userID uint) (string, error) {
	claims := newAccessJWTClaims(userID)
	return newSignedJWT(claims)
}

func (s *JWTServiceImpl) NewRefreshToken(userID uint) (string, error) {
	claims := newRefreshJWTClaims(userID)
	return newSignedJWT(claims)
}

func (s *JWTServiceImpl) ParseToken(tokenString string) (*JWTClaims, error) {
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
	return newJWTClaims(userID, "access", viper.GetDuration("jwt.access_lifetime"))
}

func newRefreshJWTClaims(userID uint) *JWTClaims {
	return newJWTClaims(userID, "refresh", viper.GetDuration("jwt.refresh_lifetime"))
}

func newJWTClaims(userID uint, tokenType string, lifetime time.Duration) *JWTClaims {
	return &JWTClaims{
		UserID: userID,
		UID:    uuid.New().String(),
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(lifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "l2sh.admissions",
		},
	}
}
