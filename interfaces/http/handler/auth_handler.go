package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/service"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

type AuthHandler interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
}

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return &authHandler{service: service}
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	logData := helpers.CreateLog(c)
	var input model.LoginInput

	if err := c.BodyParser(&input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Log:     &logData,
		})
	}

	if err := helpers.ValidateInput(input); err != nil {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  fiber.StatusBadRequest,
			Success: false,
			Message: "Invalid or malformed request body",
			Errors:  err,
			Log:     &logData,
		})
	}

	model.SanitizeLoginInput(&input)

	response := h.service.Login(c.Context(), &input)
	response.Log = &logData

	return helpers.ResponseFormatter(c, response)
}

func (h *authHandler) Logout(c *fiber.Ctx) error {
	// log := helpers.CreateLog(c)
	return nil
}

func (h *authHandler) Refresh(c *fiber.Ctx) error {
	// log := helpers.CreateLog(c)
	return nil
}

func (h *authHandler) Verify(c *fiber.Ctx) error {
	// log := helpers.CreateLog(c)
	return nil
}
