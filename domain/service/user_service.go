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
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	user, err := s.repository.FindByUUID(ctx, uuid)
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

	userModel := model.UserToDetailModel(user)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User data found",
		Data:    userModel,
	}
}

func (s *userService) GetAll(ctx context.Context, query *model.QueryGet, url string) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	users, err := s.repository.FindAll(ctx, query)
	if users == nil || err != nil {
		logData.Message = "User Not Found"
		logData.Err = err
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: logData.Message,
			Errors:  logData.Err,
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
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	userEntity := model.UserInputToEntity(input)

	if userEntity == nil {
		logData.Message = "Error parsing model"
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: logData.Message,
		}
	}

	if err := s.ValidateEntityInput(ctx, userEntity); err != nil {
		logData.Message = "invalid or malformed request body"
		logData.Err = err
		return helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: logData.Message,
			Errors:  logData.Err,
		}
	}

	if err := s.repository.Insert(ctx, userEntity); err != nil {
		logData.Message = "Error creating data"
		logData.Err = err
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: "Error creating data",
			Errors:  logData.Err,
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusCreated,
		Success: true,
		Message: "User successfully created",
	}

}

func (s *userService) UpdateByID(ctx context.Context, input *model.UserUpdateInput, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if user, err := s.repository.FindByID(ctx, id); user == nil || err != nil {
		logData.Message = "User not found"
		logData.Err = err
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: logData.Message,
			Errors:  logData.Err,
		}
	}

	userEntity := model.UserUpdateInputToEntity(input)
	if userEntity == nil {
		logData.Message = "Error parsing model"
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: logData.Message,
		}
	}

	userEntity.ID = id

	if err := s.ValidateEntityInput(ctx, userEntity); err != nil {
		logData.Message = "invalid or malformed request body"
		logData.Err = err
		return helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: logData.Message,
			Errors:  logData.Err,
		}
	}

	if err := s.repository.Update(ctx, userEntity); err != nil {
		logData.Message = "Error updating data"
		logData.Err = err
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: logData.Message,
			Errors:  logData.Err,
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User successfully updated",
	}
}

func (s *userService) ChangePassByID(ctx context.Context, input *model.ChangePasswordInput, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	if user, err := s.repository.FindByID(ctx, id); user == nil || err != nil {
		logData.Message = "User not found"
		logData.Err = err
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: logData.Message,
			Errors:  logData.Err,
		}
	}

	userEntity := model.ChangePasswordToEntity(input)
	if userEntity == nil {
		logData.Message = "Error parsing model"
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: logData.Message,
		}
	}

	userEntity.ID = id
	if err := s.repository.Update(ctx, userEntity); err != nil {
		logData.Message = "Error updating data"
		logData.Err = err
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: logData.Message,
			Errors:  logData.Err,
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User password successfully updated",
	}
}

func (s *userService) DeleteByID(ctx context.Context, id uint) helpers.BaseResponse {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

	user, err := s.repository.FindByID(ctx, id)
	if err != nil || user == nil {
		logData.Message = "User not found"
		logData.Err = err
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: logData.Message,
			Errors:  logData.Err,
		}
	}

	if err := s.repository.Delete(ctx, user); err != nil {
		logData.Message = "Error deleting data"
		logData.Err = err
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Message: logData.Message,
			Errors:  logData.Err,
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User successfully deleted",
	}
}

func (s *userService) ValidateEntityInput(ctx context.Context, user *entity.User) interface{} {
	logData := helpers.CreateLog(s)
	defer helpers.LogSystemWithDefer(ctx, &logData)

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
		logData.Message = "Validation error"
		logData.Err = errors
		return errors
	}
	return nil
}
