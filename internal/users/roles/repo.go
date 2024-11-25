package roles

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"gorm.io/gorm"
)

type RolesRepo interface {
	CreateRole(role *Role) error
	RoleExists(title string) (bool, error)
	GetRoleByTitle(title string) (*Role, error)
}

type RolesRepoImpl struct {
	storage datastore.Storage
}

func NewRolesRepo(storage datastore.Storage) RolesRepo {
	if err := storage.DB.AutoMigrate(&Role{}); err != nil {
		panic(err)
	}
	return &RolesRepoImpl{storage: storage}
}

func (r *RolesRepoImpl) CreateRole(role *Role) error {
	err := r.storage.DB.Create(role).Error
	if err != nil {
		return errors.Join(errors.New("failed to create role"), err)
	}

	return nil
}

func (r *RolesRepoImpl) RoleExists(title string) (bool, error) {
	var role Role
	err := r.storage.DB.Where("title = ?", title).First(&role).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return true, nil
}

func (r *RolesRepoImpl) GetRoleByTitle(title string) (*Role, error) {
	var role Role
	err := r.storage.DB.Where("title = ?", title).First(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}
