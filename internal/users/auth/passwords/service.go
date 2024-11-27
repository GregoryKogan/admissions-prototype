package passwords

import (
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords/crypto"
	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
)

var (
	ErrPasswordTooShort  = fmt.Errorf("password must be at least %d characters long", viper.GetInt("auth.passwords.min_length"))
	ErrPasswordNoNumber  = errors.New("password must contain at least one number")
	ErrPasswordNoUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoSpecial = errors.New("password must contain at least one special character")
	ErrFailedToGenerate  = errors.New("failed to generate password")
)

type PasswordsService interface {
	GetByUserID(userID uint) (*Password, error)
	Create(userID uint, password string) error
	Validate(password string) error
	Verify(userID uint, password string) (bool, error)
	Generate() string
}

type PasswordsServiceImpl struct {
	crypto crypto.CryptoService
	repo   PasswordsRepo
}

func NewPasswordsService(repo PasswordsRepo) PasswordsService {
	return &PasswordsServiceImpl{
		crypto: crypto.NewCryptoService(),
		repo:   repo,
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

	hashedPassword, err := s.crypto.GenerateHash([]byte(password))
	if err != nil {
		return fmt.Errorf("failed to generate hash: %w", err)
	}

	return s.repo.Create(userID, hashedPassword)
}

func (s *PasswordsServiceImpl) Validate(password string) error {
	if len(password) < viper.GetInt("auth.passwords.min_length") {
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

	hashedPassword := &crypto.HashedPassword{
		Hash:      record.Hash,
		Salt:      record.Salt,
		Algorithm: record.Algorithm,
	}

	match, err := s.crypto.Compare([]byte(password), hashedPassword)
	if err != nil {
		return false, fmt.Errorf("failed to compare passwords: %w", err)
	}

	return match, nil
}

func (s *PasswordsServiceImpl) Generate() string {
	const (
		lowerChars   = "abcdefghijklmnopqrstuvwxyz"
		upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numberChars  = "0123456789"
		specialChars = "!@#$%^&*()_+{}|:<>?~[]"
	)

	length := viper.GetInt("auth.passwords.gen_length")
	if length < 4 { // Ensure minimum length for all character types
		length = 4
	}

	seed := uint64(time.Now().UnixNano())
	rng := rand.New(rand.NewSource(seed))

	// Initialize with one character from each required set
	password := []byte{
		lowerChars[rng.Intn(len(lowerChars))],
		upperChars[rng.Intn(len(upperChars))],
		numberChars[rng.Intn(len(numberChars))],
		specialChars[rng.Intn(len(specialChars))],
	}

	// Fill the rest with random characters from all sets
	allChars := lowerChars + upperChars + numberChars + specialChars
	for i := 4; i < length; i++ {
		password = append(password, allChars[rng.Intn(len(allChars))])
	}

	// Shuffle the password
	for i := len(password) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password)
}
