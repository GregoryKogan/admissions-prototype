package users

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UsersService interface {
	AddUserPreloadMiddleware(g *echo.Group) error
	GetByEmail(email string) (*User, error)
	GetByID(userID uint) (*User, error)
	Create(email string) (*User, error)
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

func (s *UsersServiceImpl) GetByEmail(email string) (*User, error) {
	return s.repo.GetByEmail(email)
}

func (s *UsersServiceImpl) GetByID(userID uint) (*User, error) {
	return s.repo.GetByID(userID)
}

func (s *UsersServiceImpl) Create(email string) (*User, error) {
	// Check if user with the same email already exists
	user, err := s.GetByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Join(errors.New("failed to get user by email"), err)
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
		Email:  email,
		RoleID: role.ID,
		Role:   *role,
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
