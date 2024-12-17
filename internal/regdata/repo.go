package regdata

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"gorm.io/gorm"
)

type RegistrationDataRepo interface {
	Create(data *RegistrationData) error
	Delete(id uint) error
	GetByID(id uint) (*RegistrationData, error)
	ExistsByEmailNameAndGrade(email, name string, grade uint) (bool, error)
	SetEmailVerified(registrationID uint) error
	GetPending() ([]*RegistrationData, error)
	GetAccepted() ([]*RegistrationData, error)
}

type RegistrationDataRepoImpl struct {
	storage datastore.Storage
}

func NewRegistrationDataRepo(storage datastore.Storage) RegistrationDataRepo {
	if err := storage.DB().AutoMigrate(&RegistrationData{}); err != nil {
		panic(err)
	}
	return &RegistrationDataRepoImpl{storage: storage}
}

func (r *RegistrationDataRepoImpl) Create(data *RegistrationData) error {
	err := r.storage.DB().Create(data).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *RegistrationDataRepoImpl) Delete(id uint) error {
	result := r.storage.DB().Delete(&RegistrationData{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	return nil
}

func (r *RegistrationDataRepoImpl) GetByID(id uint) (*RegistrationData, error) {
	var data RegistrationData
	err := r.storage.DB().First(&data, id).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *RegistrationDataRepoImpl) ExistsByEmailNameAndGrade(email, name string, grade uint) (bool, error) {
	var data RegistrationData
	err := r.storage.DB().Where("email = ? AND first_name = ? AND grade = ?", email, name, grade).First(&data).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RegistrationDataRepoImpl) SetEmailVerified(registrationID uint) error {
	result := r.storage.DB().Model(&RegistrationData{}).Where("id = ?", registrationID).Update("email_verified", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}

	return nil
}

func (r *RegistrationDataRepoImpl) GetPending() ([]*RegistrationData, error) {
	var registrations []*RegistrationData
	// email_verified = true and no user has RegistrationDataID set to this record
	err := r.storage.DB().Model(&RegistrationData{}).
		Joins("LEFT JOIN users ON users.registration_data_id = registration_data.id").
		Where("email_verified = ? AND users.id IS NULL", true).
		Find(&registrations).Error
	if err != nil {
		return nil, err
	}
	return registrations, nil
}

func (r *RegistrationDataRepoImpl) GetAccepted() ([]*RegistrationData, error) {
	var registrations []*RegistrationData
	err := r.storage.DB().Model(&RegistrationData{}).
		Joins("JOIN users ON users.registration_data_id = registration_data.id").
		Preload("User").
		Find(&registrations).Error
	if err != nil {
		return nil, err
	}
	return registrations, nil
}
