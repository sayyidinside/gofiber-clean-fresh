package user

import "gorm.io/gorm"

type UserRepository interface {
	FindByID(id string) (*User, error)
	// Create(*User) error
}

type repository struct {
	*gorm.DB
}

func NewRepository(db *gorm.DB) UserRepository {
	return &repository{DB: db}
}

func (r *repository) FindByID(id string) (*User, error) {
	var user User
	if err := r.DB.Limit(1).Where("id = ?", id).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		}).
		Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
