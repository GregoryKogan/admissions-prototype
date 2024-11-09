package auth_test

import (
	"errors"
	"testing"

	"github.com/L2SH-Dev/admissions/internal/auth"
	"github.com/L2SH-Dev/admissions/internal/passwords"
	"github.com/L2SH-Dev/admissions/internal/secrets"
	"github.com/stretchr/testify/assert"
)

type MockPasswordsService struct {
	users map[uint]string
	err   error
}

func NewMockPasswordsService() *MockPasswordsService {
	return &MockPasswordsService{
		users: make(map[uint]string),
	}
}

func (m *MockPasswordsService) GetByUserID(userID uint) (*passwords.Password, error) {
	if m.err != nil {
		return nil, m.err
	}
	password, ok := m.users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &passwords.Password{
		UserID:    userID,
		Hash:      []byte(password),
		Salt:      []byte("salt"),
		Algorithm: "sha256",
	}, nil
}

func (m *MockPasswordsService) Validate(password string) error {
	if m.err != nil {
		return m.err
	}
	// Add validation logic if necessary
	return nil
}

func (m *MockPasswordsService) Create(userID uint, password string) error {
	if m.err != nil {
		return m.err
	}
	m.users[userID] = password
	return nil
}

func (m *MockPasswordsService) Verify(userID uint, password string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	storedPassword, ok := m.users[userID]
	if !ok {
		return false, errors.New("user not found")
	}
	return storedPassword == password, nil
}

func TestAuthService_Register(t *testing.T) {
	mockPasswordsService := NewMockPasswordsService()
	authService := auth.NewAuthService(mockPasswordsService)

	err := authService.Register(1, "password123")
	assert.NoError(t, err)
	assert.Equal(t, "password123", mockPasswordsService.users[1])
}

func TestAuthService_Login_Success(t *testing.T) {
	mockPasswordsService := NewMockPasswordsService()
	authService := auth.NewAuthService(mockPasswordsService)

	secrets.SetMockSecret("jwt_key", "secret")
	defer secrets.ClearMockSecrets()

	// Register user
	err := authService.Register(1, "password123")
	assert.NoError(t, err)

	// Successful login
	tokenPair, err := authService.Login(1, "password123")
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockPasswordsService := NewMockPasswordsService()
	authService := auth.NewAuthService(mockPasswordsService)

	// Register user
	err := authService.Register(1, "password123")
	assert.NoError(t, err)

	// Attempt login with wrong password
	tokenPair, err := authService.Login(1, "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, auth.ErrInvalidPassword, err)
	assert.Nil(t, tokenPair)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockPasswordsService := NewMockPasswordsService()
	authService := auth.NewAuthService(mockPasswordsService)

	// Attempt login without registering
	tokenPair, err := authService.Login(1, "password123")
	assert.Error(t, err)
	assert.Nil(t, tokenPair)
}

func TestAuthService_ValidatePassword(t *testing.T) {
	mockPasswordsService := NewMockPasswordsService()
	authService := auth.NewAuthService(mockPasswordsService)

	// Validate valid password
	err := authService.ValidatePassword("validpassword")
	assert.NoError(t, err)

	// Simulate validation error
	mockPasswordsService.err = errors.New("validation error")
	err = authService.ValidatePassword("invalidpassword")
	assert.Error(t, err)
	assert.Equal(t, "validation error", err.Error())
}
