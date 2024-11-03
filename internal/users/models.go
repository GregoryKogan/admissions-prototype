package users

import (
	"github.com/jackc/pgx/pgtype"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string   `json:"email" gorm:"unique;not null"`
	RoleID   uint     `json:"-"`
	Role     Role     `json:"role" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Password Password `json:"-"`
}

type Password struct {
	gorm.Model
	UserID    uint   `gorm:"index;unique;not null"`
	Hash      string `gorm:"not null"`
	Salt      string `gorm:"unique;not null"`
	Algorithm string `gorm:"not null"`
}

type Role struct {
	gorm.Model
	Title       string       `json:"title" gorm:"index;unique;not null"`
	Permissions pgtype.JSONB `json:"permissions" gorm:"type:jsonb;not null"`
}
