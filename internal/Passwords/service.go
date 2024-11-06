package passwords

import (
	"errors"
	"fmt"
)

type PasswordsService struct {
	repo *PasswordsRepo
}

func NewPasswordsService(repo *PasswordsRepo) *PasswordsService {
	return &PasswordsService{repo: repo}
}

func (s *PasswordsService) GetByUserID(userID uint) (*Password, error) {
	record, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, errors.Join(errors.New("failed to get password record"), err)
	}

	return record, nil
}

func (s *PasswordsService) Create(userID uint, password string) error {
	exists, err := s.repo.ExistsByUserID(userID)
	if err != nil {
		return errors.Join(errors.New("failed to check if password exists"), err)
	}

	if exists {
		return fmt.Errorf("password already exists for user with ID: %d", userID)
	}

	params := DefaultArgon2idParams()
	hashedPassword, err := params.GenerateHash([]byte(password), nil)
	if err != nil {
		return err
	}

	return s.repo.Create(userID, hashedPassword)
}

func (s *PasswordsService) Verify(userID uint, password string) (bool, error) {
	record, err := s.GetByUserID(userID)
	if err != nil {
		return false, err
	}

	hashedPassword := &HashedPassword{
		Hash:      []byte(record.Hash),
		Salt:      []byte(record.Salt),
		Algorithm: record.Algorithm,
	}

	params := DefaultArgon2idParams()
	same, err := params.Compare([]byte(password), hashedPassword)
	if err != nil {
		return false, errors.Join(errors.New("failed to compare passwords"), err)
	}

	return same, nil
}
