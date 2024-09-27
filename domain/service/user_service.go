package service

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uint) helpers.BaseResponse
	Create()
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		repository: repository,
	}
}

func (s *userService) GetUserByID(ctx context.Context, id uint) helpers.BaseResponse {
	user, err := s.repository.FindByID(id)
	if user == nil || err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusNotFound,
			Success: false,
			Message: "User Not Found",
		}
	}

	// convert entity to model data
	userModel := s.entityToDetailModel(user)
	iUserModel := interface{}(userModel)

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: "User data found",
		Data:    &iUserModel,
	}
}

func (s *userService) Create() {
	log.Println("test")
}

func (s *userService) entityToDetailModel(user *entity.User) *model.UserDetail {
	return &model.UserDetail{
		ID:          user.ID,
		UUID:        user.UUID,
		RoleID:      user.RoleID,
		Role:        user.Role.Name,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		ValidatedAt: user.ValidatedAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
