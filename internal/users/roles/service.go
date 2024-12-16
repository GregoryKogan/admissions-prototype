package roles

import (
	"errors"

	"github.com/spf13/viper"
)

type RolesService interface {
	CreateRole(role *Role) error
	CreateDefaultRoles() error
	RoleExists(title string) (bool, error)
	GetRoleByTitle(title string) (*Role, error)
}

type RolesServiceImpl struct {
	repo RolesRepo
}

func NewRolesService(repo RolesRepo) RolesService {
	return &RolesServiceImpl{repo: repo}
}

func (s *RolesServiceImpl) CreateRole(role *Role) error {
	if exists, err := s.RoleExists(role.Title); err != nil {
		return errors.Join(errors.New("failed to check if role exists"), err)
	} else if exists {
		return errors.New("role already exists")
	}

	return s.repo.CreateRole(role)
}

func (s *RolesServiceImpl) CreateDefaultRoles() error {
	rolesConfig := viper.GetStringMap("users.roles")

	for roleTitle, roleData := range rolesConfig {
		permissions := roleData.(map[string]interface{})["permissions"].(map[string]interface{})

		role := Role{
			Title:        roleTitle,
			Admin:        permissions["admin"].(bool),
			WriteGeneral: permissions["write_general"].(bool),
			AIAccess:     permissions["ai_access"].(bool),
		}

		if exists, err := s.RoleExists(role.Title); err != nil {
			return errors.Join(errors.New("failed to check if role exists"), err)
		} else if exists {
			continue
		}

		err := s.repo.CreateRole(&role)
		if err != nil {
			return errors.Join(errors.New("failed to create role"), err)
		}
	}

	return nil
}

func (s *RolesServiceImpl) RoleExists(title string) (bool, error) {
	return s.repo.RoleExists(title)
}

func (s *RolesServiceImpl) GetRoleByTitle(title string) (*Role, error) {
	return s.repo.GetRoleByTitle(title)
}
