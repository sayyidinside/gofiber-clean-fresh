package database

import (
	"github.com/sayyidinside/gofiber-clean-fresh/domain/user"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&user.User{})
}
