package users

import (
	"encoding/json"

	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/jackc/pgx/pgtype"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string             `json:"email" gorm:"unique;not null"`
	RoleID   uint               `json:"-"`
	Role     Role               `json:"role" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Password passwords.Password `json:"-" gorm:"not null"`
}

type Role struct {
	gorm.Model
	Title       string       `json:"title" gorm:"index;unique;not null"`
	Permissions pgtype.JSONB `json:"permissions" gorm:"type:jsonb;not null"`
}

func (r Role) MarshalJSON() ([]byte, error) {
	// This is a custom implementation of the MarshalJSON method for the Role struct.
	// We need to convert the pgtype.JSONB field to a map[string]interface{} before marshaling the struct.

	type Alias Role
	permissions := map[string]interface{}{}

	if err := json.Unmarshal(r.Permissions.Bytes, &permissions); err != nil {
		return nil, err
	}
	return json.Marshal(&struct {
		Permissions map[string]interface{} `json:"permissions"`
		Alias
	}{
		Permissions: permissions,
		Alias:       (Alias)(r),
	})
}
