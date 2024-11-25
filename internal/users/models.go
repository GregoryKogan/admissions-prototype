package users

import (
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string             `json:"email" gorm:"unique;not null"`
	RoleID   uint               `json:"-"`
	Role     roles.Role         `json:"role" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Password passwords.Password `json:"-" gorm:"not null"`
}
