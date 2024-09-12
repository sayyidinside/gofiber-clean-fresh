package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handlers"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes/v1/users"
)

func RegisterRoutes(route fiber.Router, handler *handlers.Handlers) {
	v1 := route.Group("/v1")

	users.RegisterRoutes(v1, handler.UserManagementHandler)
}
