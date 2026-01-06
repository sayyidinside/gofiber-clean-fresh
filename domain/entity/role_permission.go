package entity

import "github.com/sayyidinside/gofiber-clean-fresh/pkg/utils/constant"

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"`
	PermissionID uint `gorm:"primaryKey"`

	// Relationships
	Role       Role       `gorm:"foreignKey:RoleID"`
	Permission Permission `gorm:"foreignKey:PermissionID"`
}

func (RolePermission) TableName() string {
	return constant.TABLE_ROLE_PERMISSION
}
