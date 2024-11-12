package auth

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/users/auth/authjwt"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/labstack/echo/v4"
)

var ErrInvalidToken = errors.New("invalid token")

type AuthService interface {
	AddAuthMiddleware(g *echo.Group) error
	ValidatePassword(password string) error
	IsTokenCached(claims *authjwt.JWTClaims) (bool, error)
	Register(userID uint, password string) error
	Login(userID uint, password string) (*TokenPair, error)
	Refresh(refreshToken string) (*TokenPair, error)
	Logout(userID uint)
}

type AuthServiceImpl struct {
	repo             AuthRepo
	passwordsService passwords.PasswordsService
	jwtService       authjwt.JWTService
}

func NewAuthService(repo AuthRepo, passwordsService passwords.PasswordsService) AuthService {
	return &AuthServiceImpl{
		repo:             repo,
		passwordsService: passwordsService,
		jwtService:       authjwt.NewJWTService(),
	}
}

var ErrInvalidPassword = errors.New("invalid password")

func (s *AuthServiceImpl) ValidatePassword(password string) error {
	return s.passwordsService.Validate(password)
}

func (s *AuthServiceImpl) IsTokenCached(claims *authjwt.JWTClaims) (bool, error) {
	return s.repo.IsTokenCached(claims)
}

func (s *AuthServiceImpl) Register(userID uint, password string) error {
	return s.passwordsService.Create(userID, password)
}

func (s *AuthServiceImpl) Login(userID uint, password string) (*TokenPair, error) {
	ok, err := s.passwordsService.Verify(userID, password)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrInvalidPassword
	}

	return s.generateTokenPair(userID)
}

func (s *AuthServiceImpl) Refresh(refreshToken string) (*TokenPair, error) {
	claims, err := s.jwtService.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, errors.Join(ErrInvalidToken, errors.New("invalid token type"))
	}

	ok, err := s.IsTokenCached(claims)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.Join(ErrInvalidToken, errors.New("token not found"))
	}

	return s.generateTokenPair(claims.UserID)
}

func (s *AuthServiceImpl) Logout(userID uint) {
	s.repo.DeleteTokenPair(userID)
}

func (s *AuthServiceImpl) generateTokenPair(userID uint) (*TokenPair, error) {
	accessToken, err := s.jwtService.NewAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.NewRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	tokenPair := &TokenPair{
		Access:  accessToken,
		Refresh: refreshToken,
	}

	err = s.repo.CacheTokenPair(tokenPair)
	if err != nil {
		return nil, errors.Join(errors.New("failed to cache token pair"), err)
	}

	return tokenPair, nil
}
