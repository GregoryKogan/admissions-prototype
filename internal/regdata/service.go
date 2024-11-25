package regdata

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/validation"
)

var (
	ErrRegistrationDataInvalid  = errors.New("registration data is invalid")
	ErrorRegistrationDataExists = errors.New("registration data with the same email, first name, and grade already exists")
)

type RegistrationDataService interface {
	CreateRegistrationData(data *RegistrationData) error
	GetByID(id uint) (*RegistrationData, error)
	SetEmailVerified(registrationID uint) error
}

type RegistrationDataServiceImpl struct {
	repo RegistrationDataRepo
}

func NewRegistrationDataService(repo RegistrationDataRepo) RegistrationDataService {
	return &RegistrationDataServiceImpl{repo: repo}
}

func (s *RegistrationDataServiceImpl) CreateRegistrationData(data *RegistrationData) error {
	validator := validation.NewCustomValidator()
	err := validator.Validate(data)
	if err != nil {
		return errors.Join(ErrRegistrationDataInvalid, err)
	}

	exists, err := s.existsByEmailNameAndGrade(data.Email, data.FirstName, data.Grade)
	if err != nil {
		return err
	}
	if exists {
		return ErrorRegistrationDataExists
	}

	return s.repo.CreateRegistrationData(data)
}

func (s *RegistrationDataServiceImpl) GetByID(id uint) (*RegistrationData, error) {
	return s.repo.GetByID(id)
}

func (s *RegistrationDataServiceImpl) SetEmailVerified(registrationID uint) error {
	return s.repo.SetEmailVerified(registrationID)
}

func (s *RegistrationDataServiceImpl) existsByEmailNameAndGrade(email, name string, grade uint) (bool, error) {
	return s.repo.ExistsByEmailNameAndGrade(email, name, grade)
}
