package passwords

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type PasswordsRepo struct {
	db *gorm.DB
}

func NewPasswordsRepo(db *gorm.DB) *PasswordsRepo {
	if err := db.AutoMigrate(&Password{}); err != nil {
		panic(err)
	}
	return &PasswordsRepo{db: db}
}

func (r *PasswordsRepo) Create(userID uint, hashedPassword *HashedPassword) error {
	record := Password{
		UserID:    userID,
		Hash:      hashedPassword.Hash,
		Salt:      hashedPassword.Salt,
		Algorithm: hashedPassword.Algorithm,
	}

	if err := r.db.Create(&record).Error; err != nil {
		return fmt.Errorf("failed to create password record: %w", err)
	}

	return nil
}

func (r *PasswordsRepo) GetByUserID(userID uint) (*Password, error) {
	var record Password
	if err := r.db.Where("user_id = ?", userID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("password not found for user ID %d", userID)
		}
		return nil, fmt.Errorf("failed to get password for user ID %d: %w", userID, err)
	}

	return &record, nil
}

func (r *PasswordsRepo) ExistsByUserID(userID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&Password{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check existence for user ID %d: %w", userID, err)
	}
	return count > 0, nil
}
