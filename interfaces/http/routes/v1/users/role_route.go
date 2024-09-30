package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
)

func RegisterRoleRoutes(route fiber.Router, handler handler.RoleHandler) {
	user := route.Group("/roles")

	user.Get("/:id", handler.GetRole)
	user.Get("/", handler.GetAllRole)

	user.Post("", handler.CreateRole)

	user.Put("/:id", handler.UpdateRole)

	user.Delete("/:id", handler.DeleteRole)
}
