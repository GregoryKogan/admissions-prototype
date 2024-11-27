package regdata

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"gorm.io/gorm"
)

type RegistrationDataRepo interface {
	CreateRegistrationData(data *RegistrationData) error
	GetByID(id uint) (*RegistrationData, error)
	ExistsByEmailNameAndGrade(email, name string, grade uint) (bool, error)
	SetEmailVerified(registrationID uint) error
	GetAll() ([]*RegistrationData, error)
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

func (r *RegistrationDataRepoImpl) CreateRegistrationData(data *RegistrationData) error {
	err := r.storage.DB().Create(data).Error
	if err != nil {
		return err
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

func (r *RegistrationDataRepoImpl) GetAll() ([]*RegistrationData, error) {
	var registrations []*RegistrationData
	err := r.storage.DB().Find(&registrations).Error
	if err != nil {
		return nil, err
	}
	return registrations, nil
}
