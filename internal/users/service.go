package users

import (
	"errors"

	"github.com/jackc/pgx/pgtype"
)

type UsersService struct {
	repo *UsersRepo
}

func NewUsersService(repo *UsersRepo) *UsersService {
	service := &UsersService{repo: repo}
	err := service.createDefaultRoles()
	if err != nil {
		panic(err)
	}

	return service
}

func (s *UsersService) GetByEmail(email string) (*User, error) {
	return s.repo.GetByEmail(email)
}

func (s *UsersService) createDefaultRoles() error {
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
		exists, err := s.repo.RoleExists(role.Title)
		if err != nil {
			return errors.Join(errors.New("failed to check if role exists"), err)
		}

		if exists {
			continue
		}

		err = s.repo.CreateRole(&role)
		if err != nil {
			return errors.Join(errors.New("failed to create role"), err)
		}
	}

	return nil
}
