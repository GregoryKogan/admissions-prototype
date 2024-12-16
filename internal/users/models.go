package users

import (
	"github.com/L2SH-Dev/admissions/internal/users/auth/passwords"
	"github.com/L2SH-Dev/admissions/internal/users/roles"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login              string             `json:"login" gorm:"not null;unique"`
	RegistrationDataID uint               `json:"registration_id" gorm:"unique;index;not null"`
	RoleID             uint               `json:"-" gorm:"not null"`
	Role               roles.Role         `json:"role" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Password           passwords.Password `json:"-" gorm:"not null; constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (u *User) BeforeDelete(tx *gorm.DB) error {
	if u.ID == 0 {
		return nil
	}

	if u.Password.ID != 0 {
		if err := tx.Delete(&u.Password).Error; err != nil {
			return err
		}
	}

	// delete registration data
	// cant user RegistrationData here because of circular dependency
	if err := tx.Exec("DELETE FROM registration_data WHERE user_id = ?", u.ID).Error; err != nil {
		return err
	}

	// delete exam registrations
	// cant user ExamRegistration here because of circular dependency
	if err := tx.Exec("DELETE FROM exam_registrations WHERE user_id = ?", u.ID).Error; err != nil {
		return err
	}

	// delete exam results
	// cant user ExamResult here because of circular dependency
	if err := tx.Exec("DELETE FROM exam_results WHERE user_id = ?", u.ID).Error; err != nil {
		return err
	}

	return nil
}
