package users

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"gorm.io/gorm"
)

type UsersService interface {
	GetByID(userID uint) (*User, error)
	GetByLogin(login string) (*User, error)
	Create(registrationID uint, login string) (*User, error)
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
	user, err := s.getByRegistrationID(registrationID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Join(errors.New("failed to get user by registration id"), err)
	}
	if user != nil {
		return nil, ErrUserAlreadyExists
	}

	// Get default role
	role, err := s.rolesService.GetRoleByTitle("user")
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
	err = s.repo.CreateUser(newUser)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create user"), err)
	}

	return newUser, nil
}

func (s *UsersServiceImpl) Delete(userID uint) error {
	exists, err := s.repo.UserExistsByID(userID)
	if err != nil {
		return errors.Join(errors.New("failed to check if user exists"), err)
	}
	if exists {
		return s.repo.DeleteUser(userID)
	}

	return nil
}

func (s *UsersServiceImpl) getByRegistrationID(registrationID uint) (*User, error) {
	return s.repo.GetByRegistrationID(registrationID)
}
