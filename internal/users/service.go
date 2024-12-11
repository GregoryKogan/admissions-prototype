package users

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type UsersService interface {
	GetByID(userID uint) (*User, error)
	GetByLogin(login string) (*User, error)
	Create(registrationID uint, login string) (*User, error)
	CreateDefaultAdmin(registrationID uint) (*User, error)
	Delete(userID uint) error
}

type UsersServiceImpl struct {
	repo         UsersRepo
	rolesService roles.RolesService
}

var ErrUserAlreadyExists = errors.New("user with the same email already exists")

func NewUsersService(repo UsersRepo, rolesService roles.RolesService) UsersService {
	service := &UsersServiceImpl{repo: repo, rolesService: rolesService}
	err := rolesService.CreateDefaultRoles()
	if err != nil {
		panic(err)
	}

	return service
}

func (s *UsersServiceImpl) GetByID(userID uint) (*User, error) {
	return s.repo.GetByID(userID)
}

func (s *UsersServiceImpl) GetByLogin(login string) (*User, error) {
	return s.repo.GetByLogin(login)
}

func (s *UsersServiceImpl) Create(registrationID uint, login string) (*User, error) {
	// Check if user with the same registration id already exists
	user, err := s.repo.GetByRegistrationID(registrationID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Join(errors.New("failed to get user by registration id"), err)
	}
	if user != nil {
		return nil, ErrUserAlreadyExists
	}

	// Get default role
	role, err := s.rolesService.GetRoleByTitle(viper.GetString("users.default_role"))
	if err != nil {
		return nil, errors.Join(errors.New("failed to get default role"), err)
	}

	// Create user object
	newUser := &User{
		Login:              login,
		RegistrationDataID: registrationID,
		RoleID:             role.ID,
	}

	// Create user
	err = s.repo.Create(newUser)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create user"), err)
	}

	// Verify that the user was created successfully
	createdUser, err := s.repo.GetByID(newUser.ID)
	if err != nil {
		return nil, errors.Join(errors.New("failed to retrieve created user"), err)
	}

	return createdUser, nil
}

func (s *UsersServiceImpl) CreateDefaultAdmin(registrationID uint) (*User, error) {
	// Check if user with the same registration id already exists
	user, err := s.repo.GetByRegistrationID(registrationID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Join(errors.New("failed to get user by registration id"), err)
	}
	if user != nil {
		return nil, ErrUserAlreadyExists
	}

	// Get admin role
	role, err := s.rolesService.GetRoleByTitle(viper.GetString("users.default_admin.role"))
	if err != nil {
		return nil, errors.Join(errors.New("failed to get default admin role"), err)
	}

	// Create user object
	newUser := &User{
		Login:              viper.GetString("users.default_admin.login"),
		RegistrationDataID: registrationID,
		RoleID:             role.ID,
	}

	// Create user
	err = s.repo.Create(newUser)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create user"), err)
	}

	// Verify that the user was created successfully
	createdUser, err := s.repo.GetByID(newUser.ID)
	if err != nil {
		return nil, errors.Join(errors.New("failed to retrieve created user"), err)
	}

	return createdUser, nil
}

func (s *UsersServiceImpl) Delete(userID uint) error {
	exists, err := s.repo.ExistsByID(userID)
	if err != nil {
		return errors.Join(errors.New("failed to check if user exists"), err)
	}
	if exists {
		return s.repo.Delete(userID)
	}

	return nil
}
