package passwords

import (
	"errors"

	"gorm.io/gorm"
)

type PasswordsRepo struct {
	db *gorm.DB
}

func NewPasswordRepo(db *gorm.DB) *PasswordsRepo {
	if err := db.AutoMigrate(&Password{}); err != nil {
		panic(err)
	}
	return &PasswordsRepo{db: db}
}

func (r *PasswordsRepo) Create(userID uint, hashedPassword *HashedPassword) error {
	record := Password{
		UserID:    userID,
		Hash:      string(hashedPassword.Hash),
		Salt:      string(hashedPassword.Salt),
		Algorithm: hashedPassword.Algorithm,
	}

	err := r.db.Create(&record).Error
	if err != nil {
		return errors.Join(errors.New("failed to create password record"), err)
	}

	return nil
}

func (r *PasswordsRepo) GetByUserID(userID uint) (*Password, error) {
	var record Password
	err := r.db.Where("user_id = ?", userID).First(&record).Error
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *PasswordsRepo) ExistsByUserID(userID uint) (bool, error) {
	var record Password
	err := r.db.Where("user_id = ?", userID).First(&record).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return true, nil
}
