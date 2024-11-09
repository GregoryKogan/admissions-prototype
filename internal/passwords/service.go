package passwords

import (
	"errors"
	"fmt"
	"unicode"
)

var (
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters long")
	ErrPasswordNoNumber  = errors.New("password must contain at least one number")
	ErrPasswordNoUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoSpecial = errors.New("password must contain at least one special character")
)

type PasswordsService interface {
	GetByUserID(userID uint) (*Password, error)
	Create(userID uint, password string) error
	Validate(password string) error
	Verify(userID uint, password string) (bool, error)
}

type PasswordsServiceImpl struct {
	repo   *PasswordsRepo
	params *Argon2idParams
}

func NewPasswordsService(repo *PasswordsRepo) *PasswordsServiceImpl {
	return &PasswordsServiceImpl{
		repo:   repo,
		params: DefaultArgon2idParams(),
	}
}

func (s *PasswordsServiceImpl) GetByUserID(userID uint) (*Password, error) {
	record, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get password record: %w", err)
	}
	return record, nil
}

func (s *PasswordsServiceImpl) Create(userID uint, password string) error {
	if err := s.Validate(password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	exists, err := s.repo.ExistsByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to check if password exists: %w", err)
	}
	if exists {
		return fmt.Errorf("password already exists for user with ID: %d", userID)
	}

	hashedPassword, err := s.params.GenerateHash([]byte(password))
	if err != nil {
		return fmt.Errorf("failed to generate hash: %w", err)
	}

	return s.repo.Create(userID, hashedPassword)
}

func (s *PasswordsServiceImpl) Validate(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	var (
		hasNumber  bool
		hasUpper   bool
		hasLower   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasNumber {
		return ErrPasswordNoNumber
	}
	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}

func (s *PasswordsServiceImpl) Verify(userID uint, password string) (bool, error) {
	record, err := s.GetByUserID(userID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve password: %w", err)
	}

	hashedPassword := &HashedPassword{
		Hash:      record.Hash,
		Salt:      record.Salt,
		Algorithm: record.Algorithm,
	}

	match, err := s.params.Compare([]byte(password), hashedPassword)
	if err != nil {
		return false, fmt.Errorf("failed to compare passwords: %w", err)
	}

	return match, nil
}
