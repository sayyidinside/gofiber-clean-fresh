package service

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type PermissionService interface {
	GetByID(ctx context.Context, id uint) helpers.BaseResponse
	GetAll(ctx context.Context, query *model.QueryGet) helpers.BaseResponse
	Create()
}

type permissionService struct {
	repository       repository.PermissionRepository
	moduleRepository repository.ModuleRepository
}

func NewPermissionService(repository repository.PermissionRepository, moduleRepository repository.ModuleRepository) PermissionService {
	return &permissionService{
		repository:       repository,
		moduleRepository: moduleRepository,
	}
}

func (s *permissionService) GetByID(ctx context.Context, id uint) helpers.BaseResponse {
	permission, err := s.repository.FindByID(ctx, id)
	if permission == nil || err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Permission not found",
		}
	}

	module, _ := s.moduleRepository.FindByIDUnscoped(ctx, permission.ModuleID)

	// convert entity to model data
	permissionModel := model.PermissionToDetailModel(permission, module.Name)

	data := interface{}(permissionModel)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Permission data found",
		Data:    &data,
	}
}

func (s *permissionService) GetAll(ctx context.Context, query *model.QueryGet) helpers.BaseResponse {
	permissions, err := s.repository.FindAll(ctx, query)
	if permissions == nil || err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "Permission not found",
		}
	}

	permissionModels := model.PermissionToListModels(permissions)

	data := interface{}(permissionModels)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "Permission data found",
		Data:    &data,
	}
}

func (s *permissionService) Create() {
	log.Println("test")
}
