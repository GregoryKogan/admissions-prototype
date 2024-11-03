package users

import (
	"github.com/jackc/pgx/pgtype"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &Password{}, &Role{})
	if err != nil {
		return err
	}
	return nil
}

func CreateDefaultRoles(db *gorm.DB) error {
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
		var existingRole Role
		err := db.Where("title = ?", role.Title).First(&existingRole).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == gorm.ErrRecordNotFound {
			err = db.Create(&role).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func HasPermission(role Role, permission string) bool {
	if role.Permissions.Status != pgtype.Present {
		return false
	}

	var permissions map[string]bool
	err := role.Permissions.AssignTo(&permissions)
	if err != nil {
		return false
	}

	return permissions[permission]
}
