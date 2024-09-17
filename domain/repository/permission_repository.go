package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.Permission, error)
	FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.Permission, error)
	// Create(*Permission) error
}

type permissionRepository struct {
	*gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{DB: db}
}

func (r *permissionRepository) FindByID(ctx context.Context, id uint) (*entity.Permission, error) {
	var permission entity.Permission
	if err := r.DB.Limit(1).Where("id = ?", id).Find(&permission).Error; err != nil {
		return nil, err
	}

	return &permission, nil
}

func (r *permissionRepository) FindByUUID(uuid uuid.UUID) (*entity.Permission, error) {
	var permission entity.Permission
	if err := r.DB.Limit(1).Where("uuid = ?", uuid).Find(&permission).Error; err != nil {
		return nil, err
	}

	return &permission, nil
}

func (r *permissionRepository) FindAll(ctx context.Context, query *model.QueryGet) (*[]entity.Permission, error) {
	var permissions []entity.Permission

	tx := r.DB.WithContext(ctx)

	// pagination
	tx = tx.Scopes(helpers.Paginate(query))

	{ // Apply Order
		var allowedFields = map[string]string{
			"name":    "permissions.name",
			"module":  "permissions.module_id",
			"updated": "permssions.updated_at",
			"created": "permssions.created_at",
		}

		tx = tx.Scopes(helpers.Order(query, allowedFields))
	}

	if err := tx.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return &permissions, nil
}

func (r *permissionRepository) Count(ctx context.Context) int64 {
	var total int64

	r.DB.Model(&entity.Permission{}).Count(&total)

	return total
}

func (r *permissionRepository) CountUnscoped(ctx context.Context) int64 {
	var total int64

	r.DB.Model(&entity.Permission{}).Unscoped().Count(&total)

	return total
}
