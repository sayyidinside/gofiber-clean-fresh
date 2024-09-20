package service

import (
	"context"
	"errors"
	"log"

	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*model.UserDetail, error)
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

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.UserDetail, error) {
	user, err := s.repository.FindByID(id)
	if user.ID == 0 {
		return nil, errors.New("not found")
	}

	// convert entity to model data
	userModel := s.entityToDetailModel(user)

	return userModel, err
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
