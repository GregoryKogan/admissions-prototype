package passwords

import (
	"errors"
	"fmt"

	"github.com/L2SH-Dev/admissions/internal/storage"
	"gorm.io/gorm"
)

type PasswordsRepo interface {
	Create(userID uint, hashedPassword *HashedPassword) error
	GetByUserID(userID uint) (*Password, error)
	ExistsByUserID(userID uint) (bool, error)
}

type PasswordsRepoImpl struct {
	storage storage.Storage
}

func NewPasswordsRepo(storage storage.Storage) PasswordsRepo {
	if err := storage.DB.AutoMigrate(&Password{}); err != nil {
		panic(err)
	}
	return &PasswordsRepoImpl{storage: storage}
}

func (r *PasswordsRepoImpl) Create(userID uint, hashedPassword *HashedPassword) error {
	record := Password{
		UserID:    userID,
		Hash:      hashedPassword.Hash,
		Salt:      hashedPassword.Salt,
		Algorithm: hashedPassword.Algorithm,
	}

	if err := r.storage.DB.Create(&record).Error; err != nil {
		return fmt.Errorf("failed to create password record: %w", err)
	}

	return nil
}

func (r *PasswordsRepoImpl) GetByUserID(userID uint) (*Password, error) {
	var record Password
	if err := r.storage.DB.Where("user_id = ?", userID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("password not found for user ID %d", userID)
		}
		return nil, fmt.Errorf("failed to get password for user ID %d: %w", userID, err)
	}

	return &record, nil
}

func (r *PasswordsRepoImpl) ExistsByUserID(userID uint) (bool, error) {
	var count int64
	if err := r.storage.DB.Model(&Password{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check existence for user ID %d: %w", userID, err)
	}
	return count > 0, nil
}
