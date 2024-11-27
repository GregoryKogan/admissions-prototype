package users

import (
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login              string             `json:"login" gorm:"not null;unique"`
	RegistrationDataID uint               `json:"-" gorm:"unique;index;not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RoleID             uint               `json:"-" gorm:"not null"`
	Role               roles.Role         `json:"role" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Password           passwords.Password `json:"-" gorm:"not null"`
}
