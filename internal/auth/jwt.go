package auth

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Type   string `json:"type"`
	jwt.StandardClaims
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
		// check whether error is validation error or not
		if validationErr, ok := err.(*jwt.ValidationError); ok {
			if validationErr.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.Join(errors.New("malformed access token"), err)
			} else if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.Join(errors.New("expired access token"), err)
			} else if validationErr.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.Join(errors.New("access token not valid yet"), err)
			} else {
				return nil, errors.Join(errors.New("failed to parse access token"), err)
			}
		} else {
			return nil, errors.Join(errors.New("failed to parse access token"), err)
		}
	}

	return token.Claims.(*JWTClaims), nil
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
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().Add(viper.GetDuration("jwt.access_lifetime")).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			NotBefore: jwt.TimeFunc().Unix(),
			Issuer:    "l2sh.admissions",
		},
	}
}

func newRefreshJWTClaims(userID uint) *JWTClaims {
	return &JWTClaims{
		UserID: userID,
		Type:   "refresh",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().Add(viper.GetDuration("jwt.refresh_lifetime")).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			NotBefore: jwt.TimeFunc().Unix(),
			Issuer:    "l2sh.admissions",
		},
	}
}
