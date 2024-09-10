package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	UUID     uuid.UUID `json:"uuid" gorm:"uniqueIndex;type:char(36)"`
	RoleID   uint      `json:"role_id"`
	Username string    `json:"username" gorm:"index"`
	Name     string    `json:"name"`
	Email    string    `json:"email" gorm:"index"`
	Password string    `json:"password"`
	gorm.Model
}

func (User) TableName() string {
	return "users"
}
