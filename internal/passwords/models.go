package passwords

import (
	"gorm.io/gorm"
)

type Password struct {
	gorm.Model
	UserID    uint   `gorm:"index;unique;not null"`
	Hash      []byte `gorm:"type:bytea;not null"`
	Salt      []byte `gorm:"type:bytea;not null"`
	Algorithm string `gorm:"not null"`
}
