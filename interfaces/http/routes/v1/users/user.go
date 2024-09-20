package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/user"
)

func RegisterUserRoutes(route fiber.Router, handler *user.UserHandler) {
	user := route.Group("/informations")

	user.Get("/:id", handler.GetUser)
}
