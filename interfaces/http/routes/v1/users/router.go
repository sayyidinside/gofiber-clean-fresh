package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handlers"
)

func RegisterRoutes(route fiber.Router, handler *handlers.UserManagementHandler) {
	user := route.Group("/users/")

	RegisterUserRoutes(user, &handler.UserHandler)
}
