package regdata

import (
	"errors"
	"log/slog"

	"github.com/L2SH-Dev/admissions/internal/mailing"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/validation"
)

var (
	ErrRegistrationDataInvalid = errors.New("registration data is invalid")
	ErrRegistrationDataExists  = errors.New("registration data with the same email, first name, and grade already exists")
	ErrorEmailNotVerified      = errors.New("email is not verified")
)

type RegistrationDataService interface {
	CreateRegistrationData(data *RegistrationData) error
	GetByID(id uint) (*RegistrationData, error)
	SetEmailVerified(registrationID uint) error
	Accept(id uint) (*users.User, error)
	GetAll() ([]*RegistrationData, error)
}

type RegistrationDataServiceImpl struct {
	repo             RegistrationDataRepo
	usersService     users.UsersService
	authService      auth.AuthService
	passwordsService passwords.PasswordsService
}

func NewRegistrationDataService(
	repo RegistrationDataRepo,
	usersService users.UsersService,
	authService auth.AuthService,
	passwordsService passwords.PasswordsService,
) RegistrationDataService {
	return &RegistrationDataServiceImpl{
		repo:             repo,
		usersService:     usersService,
		authService:      authService,
		passwordsService: passwordsService,
	}
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
		return ErrRegistrationDataExists
	}

	return s.repo.CreateRegistrationData(data)
}

func (s *RegistrationDataServiceImpl) GetByID(id uint) (*RegistrationData, error) {
	return s.repo.GetByID(id)
}

func (s *RegistrationDataServiceImpl) SetEmailVerified(registrationID uint) error {
	return s.repo.SetEmailVerified(registrationID)
}

func (s *RegistrationDataServiceImpl) Accept(id uint) (*users.User, error) {
	regData, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if !regData.EmailVerified {
		return nil, ErrorEmailNotVerified
	}

	login := generateLogin(regData)
	user, err := s.usersService.Create(id, login)
	if err != nil {
		return nil, err
	}

	password := s.passwordsService.Generate()
	err = s.authService.Register(user.ID, password)
	if err != nil {
		slog.Info("Failed to register user, deleting user", slog.Any("user_id", user.ID))
		delErr := s.usersService.Delete(user.ID)
		if delErr != nil {
			slog.Error("Failed to delete user", slog.Any("user_id", user.ID), slog.Any("err", err))
			return nil, delErr
		}
		return nil, err
	}

	err = mailing.SendLoginAndPassword(regData.Email, login, password)
	if err != nil {
		slog.Error("Failed to send login and password", slog.Any("email", regData.Email), slog.Any("err", err))
		return nil, err
	}

	return user, nil
}

func (s *RegistrationDataServiceImpl) GetAll() ([]*RegistrationData, error) {
	return s.repo.GetAll()
}

func (s *RegistrationDataServiceImpl) existsByEmailNameAndGrade(email, name string, grade uint) (bool, error) {
	return s.repo.ExistsByEmailNameAndGrade(email, name, grade)
}

func generateLogin(regData *RegistrationData) string {
	return regData.FirstName + " " + regData.LastName
}
