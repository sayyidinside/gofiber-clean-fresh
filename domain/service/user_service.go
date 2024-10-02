package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type UserService interface {
	GetByID(ctx context.Context, id uint) helpers.BaseResponse
	GetByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse
	GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse
	Create(ctx context.Context, input *model.UserInput) helpers.BaseResponse
	UpdateByID(ctx context.Context, input *model.UserUpdateInput, id uint) helpers.BaseResponse
	ChangePassByID(ctx context.Context, input *model.ChangePasswordInput, id uint) helpers.BaseResponse
	DeleteByID(ctx context.Context, id uint) helpers.BaseResponse
}

type userService struct {
	repository     repository.UserRepository
	roleRepository repository.RoleRepository
}

func NewUserService(repository repository.UserRepository, roleRepository repository.RoleRepository) UserService {
	return &userService{
		repository:     repository,
		roleRepository: roleRepository,
	}
}

func (s *userService) GetByID(ctx context.Context, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	user, err := s.repository.FindByID(ctx, id)
	if user == nil || err != nil {
		logData.Message = "User Not Found"
		logData.Err = err
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: logData.Message,
			Errors:  logData.Err,
		}
	}

	// convert entity to model data
	userModel := model.UserToDetailModel(user)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User data found",
		Data:    userModel,
	}
}

func (s *userService) GetByUUID(ctx context.Context, uuid uuid.UUID) helpers.BaseResponse {
	user, err := s.repository.FindByUUID(ctx, uuid)
	if user == nil || err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User data not found",
		}
	}

	userModel := model.UserToDetailModel(user)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User data found",
		Data:    userModel,
	}
}

func (s *userService) GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse {
	users, err := s.repository.FindAll(ctx, query)
	if users == nil || err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusOK,
			Success: false,
			Message: "User not found",
		}
	}

	userModels := model.UserToListModel(users)

	totalData := s.repository.Count(ctx, query)
	pagination := helpers.GeneratePaginationMetadata(query, url, totalData)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User data found",
		Data:    userModels,
		Meta: &helpers.Meta{
			Pagination: pagination,
		},
	}
}

func (s *userService) Create(ctx context.Context, input *model.UserInput) helpers.BaseResponse {
	userEntity := model.UserInputToEntity(input)

	if userEntity == nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing model",
		}
	}

	if err := s.ValidateEntityInput(ctx, userEntity); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "invalid or malformed request body",
			Errors:  err,
		}
	}

	if err := s.repository.Insert(ctx, userEntity); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error creating data",
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusCreated,
		Success: true,
		Message: "User successfully created",
	}

}

func (s *userService) UpdateByID(ctx context.Context, input *model.UserUpdateInput, id uint) helpers.BaseResponse {
	if user, err := s.repository.FindByID(ctx, id); user == nil || err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User not found",
		}
	}

	userEntity := model.UserUpdateInputToEntity(input)
	if userEntity == nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing to model",
		}
	}

	userEntity.ID = id

	if err := s.ValidateEntityInput(ctx, userEntity); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
		}
	}

	if err := s.repository.Update(ctx, userEntity); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating data",
			Errors:  err,
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User successfully updated",
	}
}

func (s *userService) ChangePassByID(ctx context.Context, input *model.ChangePasswordInput, id uint) helpers.BaseResponse {
	if user, err := s.repository.FindByID(ctx, id); user == nil || err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User not found",
		}
	}

	userEntity := model.ChangePasswordToEntity(input)
	if userEntity == nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error parsing to model",
		}
	}

	userEntity.ID = id
	if err := s.repository.Update(ctx, userEntity); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error updating data",
			Errors:  err,
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User password successfully updated",
	}
}

func (s *userService) DeleteByID(ctx context.Context, id uint) helpers.BaseResponse {
	user, err := s.repository.FindByID(ctx, id)
	if err != nil || user == nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User not found",
		}
	}

	if err := s.repository.Delete(ctx, user); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error deleting data",
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User successfully deleted",
	}
}

func (s *userService) ValidateEntityInput(ctx context.Context, user *entity.User) interface{} {
	errors := []helpers.ValidationError{}

	if role, err := s.roleRepository.FindByID(ctx, user.RoleID); role == nil || err != nil {
		errors = append(errors, helpers.ValidationError{
			Field: "role_id",
			Tag:   "not_found",
		})
	}

	if exist := s.repository.NameExist(ctx, user); exist {
		errors = append(errors, helpers.ValidationError{
			Field: "name",
			Tag:   "duplicate",
		})
	}

	if exist := s.repository.EmailExist(ctx, user); exist {
		errors = append(errors, helpers.ValidationError{
			Field: "email",
			Tag:   "duplicate",
		})
	}

	if exist := s.repository.UsernameExist(ctx, user); exist {
		errors = append(errors, helpers.ValidationError{
			Field: "username",
			Tag:   "duplicate",
		})
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}
