package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type RoleService interface {
	GetByID(ctx context.Context, id uint) helpers.BaseResponse
	GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse
	Create(ctx context.Context, input *model.RoleInput) helpers.BaseResponse
	UpdateByID(ctx context.Context, input *model.RoleInput, id uint) helpers.BaseResponse
	DeleteByID(ctx context.Context, id uint) helpers.BaseResponse
}

type roleService struct {
	repository     repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewRoleService(repository repository.RoleRepository, permissionRepo repository.PermissionRepository) RoleService {
	return &roleService{
		repository:     repository,
		permissionRepo: permissionRepo,
	}
}

func (s *roleService) GetByID(ctx context.Context, id uint) helpers.BaseResponse {
	role, err := s.repository.FindByID(ctx, id)
	if role == nil || err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Role not found",
		}
	}

	roleModel := model.RoleToDetailModel(role)

	data := interface{}(roleModel)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Role data found",
		Data:    &data,
	}
}

func (s *roleService) GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse {
	roles, err := s.repository.FindAll(ctx, query)
	if roles == nil || err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Role not found",
		}
	}

	roleModels := model.RoleToListModels(roles)

	data := interface{}(roleModels)

	totalData := s.repository.Count(ctx, query)

	pagination := helpers.GeneratePaginationMetadata(query, url, totalData)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Role data found",
		Data:    &data,
		Meta: &helpers.Meta{
			Pagination: pagination,
		},
	}
}

func (s *roleService) Create(ctx context.Context, input *model.RoleInput) helpers.BaseResponse {
	roleEntity := model.RoleInputToEntity(input)

	if err := s.validateEntityInput(ctx, roleEntity); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		}
	}

	permissions, err := s.permissionRepo.FindInID(ctx, input.Permissions)
	if err != nil {
		iError := interface{}(err)
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error retrieving permission data",
			Errors:  &iError,
		}
	}

	roleEntity.Permissions = permissions

	if err := s.repository.Insert(ctx, roleEntity); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error creating data",
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusCreated,
		Success: true,
		Message: "Role successfully created",
	}
}

func (s *roleService) UpdateByID(ctx context.Context, input *model.RoleInput, id uint) helpers.BaseResponse {
	// Start a new transaction
	tx := s.repository.BeginTransaction(ctx)

	// Check role existence
	role, err := s.repository.FindByID(ctx, id)
	if role == nil || err != nil {
		tx.Rollback() // Rollback on error
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Role not found",
		}
	}

	roleEntity := model.RoleInputToEntity(input)
	if roleEntity == nil {
		tx.Rollback() // Rollback on error
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing model",
		}
	}

	roleEntity.ID = id

	// Retrieve permissions
	permissions, err := s.permissionRepo.FindInID(ctx, input.Permissions)
	if err != nil || len(*permissions) == 0 {
		return helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Permission data not found",
			Errors:  err,
		}
	}

	// Validate the entity
	if err := s.validateEntityInput(ctx, roleEntity); err != nil {
		tx.Rollback() // Rollback on error
		return helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		}
	}

	// Update role entity in the database
	if err := s.repository.UpdateWithTransaction(ctx, tx, roleEntity); err != nil {
		tx.Rollback() // Rollback on error
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating data",
		}
	}

	// Replace permissions in the database
	if err := s.repository.ReplacePermissionsWithTransaction(ctx, tx, roleEntity, permissions); err != nil {
		tx.Rollback() // Rollback on error
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error replacing role permissions data",
		}
	}

	// Commit the transaction if all operations succeed
	tx.Commit()

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Role successfully updated",
	}
}

func (s *roleService) DeleteByID(ctx context.Context, id uint) helpers.BaseResponse {
	// Check modul existence
	role, err := s.repository.FindByID(ctx, id)
	if role == nil || err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Role not found",
		}
	}

	if err := s.repository.Delete(ctx, role); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error deleting data",
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Role successfully deleted",
	}
}

func (s *roleService) validateEntityInput(ctx context.Context, role *entity.Role) *interface{} {
	errs := []helpers.ValidationError{}

	// Check name duplication
	if exist := s.repository.NameExist(ctx, role); exist {
		errs = append(errs, helpers.ValidationError{
			Field: "name",
			Tag:   "duplicate",
		})
	}

	if len(errs) != 0 {
		intf := interface{}(errs)
		return &intf
	}

	return nil
}
