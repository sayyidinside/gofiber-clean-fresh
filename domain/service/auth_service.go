package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/repository"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, input *model.LoginInput) helpers.BaseResponse
}

type authService struct {
	refreshTokenRepository repository.RefreshTokenRepository
	userRepository         repository.UserRepository
}

func NewAuthService(refreshTokenRepository repository.RefreshTokenRepository, userRepository repository.UserRepository) AuthService {
	return &authService{
		refreshTokenRepository: refreshTokenRepository,
		userRepository:         userRepository,
	}
}

func (s *authService) Login(ctx context.Context, input *model.LoginInput) helpers.BaseResponse {
	cfg := config.AppConfig

	user, err := s.userRepository.FindByUsernameOrEmail(ctx, input.UsernameOrEmail)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid email / username or password",
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid email / username or password",
		}
	}

	accessToken, err := helpers.GenerateToken(user, cfg.JwtAccessTime, cfg.JwtAccessPrivateSecret, false)
	if err != nil {
		jsonData, _ := json.MarshalIndent(err.Error(), "", " ")
		log.Println(string(jsonData))
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed generating token",
		}
	}

	refreshToken, err := helpers.GenerateToken(user, cfg.JwtRefreshTime, cfg.JwtRefreshPrivateSecret, true)
	if err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed generating token",
		}
	}

	if err := s.refreshTokenRepository.Insert(ctx, &entity.RefreshToken{UserID: user.ID, Token: refreshToken}); err != nil {
		return helpers.BaseResponse{
			Status:  fiber.StatusInternalServerError,
			Success: false,
			Errors:  err,
			Message: "failed registering token",
		}
	}

	return helpers.BaseResponse{
		Status:  fiber.StatusOK,
		Success: true,
		Message: fmt.Sprintf("%s user successfully login", user.Name),
		Data: &model.AllToken{
			RefreshToken: refreshToken,
			AccessToken:  accessToken,
		},
	}
}
