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
	Count(ctx context.Context) int64
	CountUnscoped(ctx context.Context) int64
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

	tx := r.DB.WithContext(ctx).Model(&entity.Permission{}).
		Joins("JOIN modules on modules.id = permissions.module_id")

	{ // Apply Query Operation

		// map value for parsing user query input
		var allowedFields = map[string]string{
			"name":        "permissions.name",
			"module":      "permissions.module_id",
			"updated":     "permssions.updated_at",
			"created":     "permssions.created_at",
			"module_name": "modules.name",
		}

		// pagination
		tx = tx.Scopes(helpers.Paginate(query))

		// Order
		tx = tx.Scopes(helpers.Order(query, allowedFields))

		// Apply Filter
		tx = tx.Scopes(helpers.Filter(query, allowedFields))

		// Apply Search
		tx = tx.Scopes(helpers.Search(query, allowedFields))
	}

	if err := tx.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return &permissions, nil
}

func (r *permissionRepository) Count(ctx context.Context) int64 {
	var total int64

	r.DB.WithContext(ctx).Model(&entity.Permission{}).Count(&total)

	return total
}

func (r *permissionRepository) CountUnscoped(ctx context.Context) int64 {
	var total int64

	r.DB.WithContext(ctx).Model(&entity.Permission{}).Unscoped().Count(&total)

	return total
}
