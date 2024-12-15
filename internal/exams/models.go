package exams

import (
	"time"

	"github.com/L2SH-Dev/admissions/internal/users"
	"gorm.io/gorm"
)

type Exam struct {
	gorm.Model
	Start      time.Time `json:"start" gorm:"not null" validate:"required"`
	End        time.Time `json:"end"`
	Location   string    `json:"location" gorm:"not null" validate:"required"`
	Capacity   uint      `json:"capacity" gorm:"not null" validate:"required"`
	Grade      uint      `json:"grade" gorm:"not null" validate:"required,min=6,max=11"`
	ExamTypeID uint      `json:"type_id" gorm:"not null" validate:"required"`
	ExamType   ExamType  `json:"type" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type ExamType struct {
	gorm.Model
	Title      string `json:"title" gorm:"unique;index;not null"`
	Order      int    `json:"order" gorm:"not null"`
	Dismissing bool   `json:"dismissing" gorm:"not null"`
	HasPoints  bool   `json:"has_points" gorm:"not null"`
}

type ExamRegistration struct {
	gorm.Model
	ExamID uint       `json:"-" gorm:"not null;index"`
	Exam   Exam       `json:"exam" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID uint       `json:"-" gorm:"not null;index"`
	User   users.User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ExamResult struct {
	gorm.Model
	ExamID    uint       `json:"-" gorm:"not null;index"`
	Exam      Exam       `json:"exam" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID    uint       `json:"-" gorm:"not null;index"`
	User      users.User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Result    string     `json:"result" gorm:"not null" validate:"required,oneof=PASSED FAILED ABSENT"`
	Dismissed bool       `json:"dismissed" gorm:"not null"`
	Points    float32    `json:"points" gorm:"not null"`
	MaxPoints float32    `json:"max_points" gorm:"not null"`
}
