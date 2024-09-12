package database

import (
	"github.com/sayyidinside/gofiber-clean-fresh/domain/module"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/permission"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/role"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/user"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&module.Module{})
	db.AutoMigrate(&permission.Permission{})
	db.AutoMigrate(&role.Role{})
	db.AutoMigrate(&user.User{})
}
