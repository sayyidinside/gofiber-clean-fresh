package repository

import (
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(id string) (*entity.User, error)
	// Create(*User) error
}

type userRepository struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByID(id string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.Limit(1).Where("id = ?", id).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		}).
		Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
