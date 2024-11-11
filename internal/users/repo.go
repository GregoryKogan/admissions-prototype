package users

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/storage"
	"gorm.io/gorm"
)

type UsersRepo interface {
	CreateRole(role *Role) error
	RoleExists(title string) (bool, error)
	GetRoleByTitle(title string) (*Role, error)

	CreateUser(user *User) error
	DeleteUser(userID uint) error
	GetByID(userID uint) (*User, error)
	GetWithDetailsByID(userID uint) (*User, error)
	UserExistsByID(userID uint) (bool, error)
	GetByEmail(email string) (*User, error)
}

type UsersRepoImpl struct {
	storage storage.Storage
}

func NewUsersRepo(storage storage.Storage) UsersRepo {
	if err := storage.DB.AutoMigrate(&User{}, &Role{}); err != nil {
		panic(err)
	}
	return &UsersRepoImpl{storage: storage}
}

func (r *UsersRepoImpl) CreateRole(role *Role) error {
	err := r.storage.DB.Create(role).Error
	if err != nil {
		return errors.Join(errors.New("failed to create role"), err)
	}

	return nil
}

func (r *UsersRepoImpl) RoleExists(title string) (bool, error) {
	var role Role
	err := r.storage.DB.Where("title = ?", title).First(&role).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return true, nil
}

func (r *UsersRepoImpl) GetRoleByTitle(title string) (*Role, error) {
	var role Role
	err := r.storage.DB.Where("title = ?", title).First(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *UsersRepoImpl) CreateUser(user *User) error {
	err := r.storage.DB.Create(user).Error
	if err != nil {
		return errors.Join(errors.New("failed to create user"), err)
	}

	return nil
}

func (r *UsersRepoImpl) DeleteUser(userID uint) error {
	err := r.storage.DB.Delete(&User{}, userID).Error
	if err != nil {
		return errors.Join(errors.New("failed to delete user"), err)
	}

	return nil
}

func (r *UsersRepoImpl) GetByID(userID uint) (*User, error) {
	var user User
	err := r.storage.DB.First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepoImpl) GetWithDetailsByID(userID uint) (*User, error) {
	var user User
	err := r.storage.DB.Preload("Role").First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepoImpl) UserExistsByID(userID uint) (bool, error) {
	user, err := r.GetByID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	if user == nil {
		return false, nil
	}

	return true, nil
}

func (r *UsersRepoImpl) GetByEmail(email string) (*User, error) {
	var user User
	err := r.storage.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
