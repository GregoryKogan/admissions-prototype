package auth

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/passwords"
)

type AuthService interface {
	ValidatePassword(password string) error
	Register(userID uint, password string) error
	Login(userID uint, password string) (*TokenPair, error)
}

type AuthServiceImpl struct {
	passwordsService passwords.PasswordsService
}

func NewAuthService(passwordsService passwords.PasswordsService) AuthService {
	return &AuthServiceImpl{
		passwordsService: passwordsService,
	}
}

var ErrInvalidPassword = errors.New("invalid password")

type TokenPair struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func (s *AuthServiceImpl) ValidatePassword(password string) error {
	return s.passwordsService.Validate(password)
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

	accessToken, err := NewAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := NewRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
