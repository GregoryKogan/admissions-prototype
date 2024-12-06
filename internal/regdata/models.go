package regdata

import (
	"time"

	"github.com/L2SH-Dev/admissions/internal/users"
	"gorm.io/gorm"
)

type RegistrationData struct {
	gorm.Model
	Email            string    `json:"email" gorm:"not null" validate:"required,email"`
	EmailVerified    bool      `json:"email_verified" gorm:"not null;default:false"`
	FirstName        string    `json:"first_name" gorm:"not null" validate:"required"`
	LastName         string    `json:"last_name" gorm:"not null" validate:"required"`
	Patronymic       string    `json:"patronymic"`
	Gender           string    `json:"gender" gorm:"varchar(1)" validate:"required,oneof=M F N"`
	BirthDate        time.Time `json:"birth_date" gorm:"not null" validate:"required"`
	Grade            uint      `json:"grade" gorm:"not null" validate:"required,min=6,max=11"`
	OldSchool        string    `json:"old_school" gorm:"not null" validate:"required"`
	ParentFirstName  string    `json:"parent_first_name" gorm:"not null" validate:"required"`
	ParentLastName   string    `json:"parent_last_name" gorm:"not null" validate:"required"`
	ParentPatronymic string    `json:"parent_patronymic"`
	ParentPhone      string    `json:"parent_phone" gorm:"not null" validate:"required,e164"`
	JuneExam         bool      `json:"june_exam" gorm:"not null;default:false"`
	VMSH             bool      `json:"vmsh" gorm:"not null;default:false"`
	Source           string    `json:"source"`
	User             users.User
}
