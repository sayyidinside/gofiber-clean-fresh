package module

import (
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/permission"
	"gorm.io/gorm"
)

type Module struct {
	ID          uint                    `json:"id" gorm:"primaryKey"`
	UUID        uuid.UUID               `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	Name        string                  `json:"name"`
	Permissions []permission.Permission `json:"permissions" gorm:"foreignKey:ModuleID"`
	gorm.Model
}

func (Module) TableName() string {
	return "modules"
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database.
// It generates a new UUID for the UUID field of the VEN_Legal struct.
func (m *Module) BeforeCreate(tx *gorm.DB) (err error) {
	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}
	return
}
