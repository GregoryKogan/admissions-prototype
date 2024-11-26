package roles

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Title        string `json:"title" gorm:"index;unique;not null"`
	Admin        bool   `json:"admin" gorm:"default:false"`
	WriteGeneral bool   `json:"write_general" gorm:"default:false"`
	AIAccess     bool   `json:"ai_access" gorm:"default:false"`
}
