package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/user"
)

type UserHandler struct {
	service user.UserService
}

func NewHandler(service user.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := uh.service.GetUserByID(id)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}
