package regdata

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/L2SH-Dev/admissions/internal/mailing"
	"github.com/L2SH-Dev/admissions/internal/users"
	"github.com/L2SH-Dev/admissions/internal/users/auth"
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/validation"
	"github.com/essentialkaos/translit/v3"
)

var (
	ErrRegistrationDataInvalid = errors.New("registration data is invalid")
	ErrRegistrationDataExists  = errors.New("registration data with the same email, first name, and grade already exists")
	ErrorEmailNotVerified      = errors.New("email is not verified")
)

type RegistrationDataService interface {
	Create(data *RegistrationData) error
	GetByID(id uint) (*RegistrationData, error)
	SetEmailVerified(registrationID uint) error
	Accept(id uint) (*users.User, error)
	Reject(id uint) error
	GetPending() ([]*RegistrationData, error)
	GetAccepted() ([]*RegistrationData, error)
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

func (s *RegistrationDataServiceImpl) Create(data *RegistrationData) error {
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

	return s.repo.Create(data)
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

func (s *RegistrationDataServiceImpl) Reject(id uint) error {
	return s.repo.Delete(id)
}

func (s *RegistrationDataServiceImpl) GetPending() ([]*RegistrationData, error) {
	return s.repo.GetPending()
}

func (s *RegistrationDataServiceImpl) GetAccepted() ([]*RegistrationData, error) {
	return s.repo.GetAccepted()
}

func (s *RegistrationDataServiceImpl) existsByEmailNameAndGrade(email, name string, grade uint) (bool, error) {
	return s.repo.ExistsByEmailNameAndGrade(email, name, grade)
}

func generateLogin(regData *RegistrationData) string {
	transliterator := translit.ICAO
	tFirstName := transliterator(strings.ToLower(regData.FirstName))[:1]
	tLastName := transliterator(strings.ToLower(regData.LastName))

	regIdStr := "00000" + fmt.Sprint(regData.ID)
	regIdStr = regIdStr[len(regIdStr)-5:]

	return fmt.Sprintf("%s.%s-%s", tFirstName, tLastName, regIdStr)
}
