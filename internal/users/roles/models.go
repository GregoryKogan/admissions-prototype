package roles

import (
	"encoding/json"

	"github.com/jackc/pgx/pgtype"
	"gorm.io/gorm"
)

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
