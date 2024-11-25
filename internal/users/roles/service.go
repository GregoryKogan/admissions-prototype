package roles

import (
	"errors"

	"github.com/jackc/pgx/pgtype"
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
	roles := []Role{
		{
			Title: "admin",
			Permissions: pgtype.JSONB{
				Bytes:  []byte(`{"admin": true}`),
				Status: pgtype.Present,
			},
		},
		{
			Title: "user",
			Permissions: pgtype.JSONB{
				Bytes:  []byte(`{"admin": false}`),
				Status: pgtype.Present,
			},
		},
	}

	for _, role := range roles {
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
