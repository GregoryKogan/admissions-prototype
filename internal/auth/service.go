package auth

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/passwords"
)

type AuthService struct {
	passwordsService *passwords.PasswordsService
}

func NewAuthService(passwordsService *passwords.PasswordsService) *AuthService {
	return &AuthService{
		passwordsService: passwordsService,
	}
}

var ErrInvalidPassword = errors.New("invalid password")

type TokenPair struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func (s *AuthService) ValidatePassword(password string) error {
	return s.passwordsService.Validate(password)
}

func (s *AuthService) Register(userID uint, password string) error {
	return s.passwordsService.Create(userID, password)
}

func (s *AuthService) Login(userID uint, password string) (*TokenPair, error) {
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
