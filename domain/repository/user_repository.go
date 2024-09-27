package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.User, error)
	FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.User, error)
	Count(ctx context.Context, query *model.QueryGet) int64
	// Create(*User) error
}

type userRepository struct {
	*gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if result := r.DB.WithContext(ctx).Limit(1).Where("id = ?", id).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		}).
		Find(&user); result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.User, error) {
	var user entity.User
	if result := r.DB.WithContext(ctx).Limit(1).Where("uuid = ?", uuid).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		}).
		Find(&user); result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.User, error) {
	var users []entity.User
	tx := r.DB.WithContext(ctx).Model(&entity.User{}).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		})

	var allowedFields = map[string]string{
		"name":    "users.name",
		"created": "users.created_at",
		"updated": "users.updated_at",
	}

	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	if err := tx.Find(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *userRepository) Count(ctx context.Context, query *model.QueryGet) int64 {
	var total int64

	tx := r.DB.WithContext(ctx).Model(&entity.User{}).
		Preload("Role", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").Unscoped()
		})

	var allowedFields = map[string]string{
		"name":    "users.name",
		"created": "users.created_at",
		"updated": "users.updated_at",
	}

	tx = tx.Scopes(
		helpers.Paginate(query),
		helpers.Order(query, allowedFields),
		helpers.Filter(query, allowedFields),
		helpers.Search(query, allowedFields),
	)

	tx.Count(&total)

	return total
}
