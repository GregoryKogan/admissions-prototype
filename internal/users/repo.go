package users

import (
	"errors"

	"gorm.io/gorm"
)

type UsersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) *UsersRepo {
	if err := db.AutoMigrate(&User{}, &Role{}); err != nil {
		panic(err)
	}
	return &UsersRepo{db: db}
}

func (r *UsersRepo) CreateRole(role *Role) error {
	err := r.db.Create(role).Error
	if err != nil {
		return errors.Join(errors.New("failed to create role"), err)
	}

	return nil
}

func (r *UsersRepo) RoleExists(title string) (bool, error) {
	var role Role
	err := r.db.Where("title = ?", title).First(&role).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return true, nil
}

func (r *UsersRepo) CreateUser(user *User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return errors.Join(errors.New("failed to create user"), err)
	}

	return nil
}

func (r *UsersRepo) GetByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
