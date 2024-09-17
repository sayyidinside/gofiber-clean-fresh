package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
)

type (
	PermissionDetail struct {
		ID        uint      `json:"id"`
		UUID      uuid.UUID `json:"uuid"`
		Name      string    `json:"name"`
		Module    string    `json:"module"`
		ModuleID  uint      `json:"module_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	PermissionList struct {
		ID       uint      `json:"id"`
		UUID     uuid.UUID `json:"uuid"`
		Name     string    `json:"name"`
		Module   string    `json:"module"`
		ModuleID uint      `json:"module_id"`
	}
)

func PermissionToDetailModel(permission *entity.Permission, moduleName string) *PermissionDetail {
	return &PermissionDetail{
		ID:        permission.ID,
		UUID:      permission.UUID,
		Name:      permission.Name,
		Module:    moduleName,
		ModuleID:  permission.ID,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}
}

func PermissionToListModel(permission *entity.Permission) *PermissionList {
	return &PermissionList{
		ID:       permission.ID,
		UUID:     permission.UUID,
		Name:     permission.Name,
		ModuleID: permission.ID,
	}
}

func PermissionToListModels(permissions *[]entity.Permission) *[]PermissionList {
	listModels := []PermissionList{}

	for _, permission := range *permissions {
		listModels = append(listModels, *PermissionToListModel(&permission))
	}

	return &listModels
}
