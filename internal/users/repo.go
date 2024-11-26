package users

import (
	"errors"

	"github.com/L2SH-Dev/admissions/internal/datastore"
	"gorm.io/gorm"
)

type UsersRepo interface {
	CreateUser(user *User) error
	DeleteUser(userID uint) error
	GetByID(userID uint) (*User, error)
	GetByRegistrationID(registrationID uint) (*User, error)
	GetByLogin(login string) (*User, error)
	UserExistsByID(userID uint) (bool, error)
}

type UsersRepoImpl struct {
	storage datastore.Storage
}

func NewUsersRepo(storage datastore.Storage) UsersRepo {
	if err := storage.DB().AutoMigrate(&User{}); err != nil {
		panic(err)
	}
	return &UsersRepoImpl{storage: storage}
}

func (r *UsersRepoImpl) CreateUser(user *User) error {
	err := r.storage.DB().Create(user).Error
	if err != nil {
		return errors.Join(errors.New("failed to create user"), err)
	}

	return nil
}

func (r *UsersRepoImpl) DeleteUser(userID uint) error {
	err := r.storage.DB().Delete(&User{}, userID).Error
	if err != nil {
		return errors.Join(errors.New("failed to delete user"), err)
	}

	return nil
}

func (r *UsersRepoImpl) GetByID(userID uint) (*User, error) {
	var user User
	err := r.storage.DB().Preload("Role").Preload("RegistrationData").First(&user, userID).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepoImpl) GetByRegistrationID(registrationID uint) (*User, error) {
	var user User
	err := r.storage.DB().Where("registration_data_id = ?", registrationID).Preload("Role").Preload("RegistrationData").First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepoImpl) GetByLogin(login string) (*User, error) {
	var user User
	err := r.storage.DB().Where("login = ?", login).Preload("Role").Preload("RegistrationData").First(&user).Error
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
