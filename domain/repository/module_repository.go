package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"gorm.io/gorm"
)

type ModuleRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Module, error)
	FindByIDUnscoped(ctx context.Context, id uint) (*entity.Module, error)
	// Create(*Module) error
}

type moduleRepository struct {
	*gorm.DB
}

func NewModuleRepository(db *gorm.DB) ModuleRepository {
	return &moduleRepository{DB: db}
}

func (r *moduleRepository) FindByID(ctx context.Context, id uint) (*entity.Module, error) {
	var module entity.Module
	if result := r.DB.WithContext(ctx).Limit(1).Where("id = ?", id).Find(&module); result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}

	return &module, nil
}

func (r *moduleRepository) FindByIDUnscoped(ctx context.Context, id uint) (*entity.Module, error) {
	var module entity.Module
	if err := r.DB.WithContext(ctx).Limit(1).Where("id = ?", id).Unscoped().Find(&module).Error; err != nil {
		return nil, err
	}

	return &module, nil
}

func (r *moduleRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*entity.Module, error) {
	var module entity.Module
	if err := r.DB.Limit(1).Where("uuid = ?", uuid).Find(&module).Error; err != nil {
		return nil, err
	}

	return &module, nil
}
