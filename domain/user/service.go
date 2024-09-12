package user

import (
	"errors"

	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
)

type UserService interface {
	GetUserByID(id string) (*model.UserDetail, error)
}

type service struct {
	repository UserRepository
}

func NewService(repository UserRepository) UserService {
	return &service{
		repository: repository,
	}
}

func (s *service) GetUserByID(id string) (*model.UserDetail, error) {
	user, err := s.repository.FindByID(id)
	if user.ID == 0 {
		return nil, errors.New("not found")
	}

	// convert entity to model data
	userModel := s.entityToDetailModel(user)

	return userModel, err
}

func (s *service) entityToDetailModel(user *User) *model.UserDetail {
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
