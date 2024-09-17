package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
)

func RegisterPermissionRoutes(route fiber.Router, handler *handler.PermissionHandler) {
	user := route.Group("/permissions")

	user.Get("/", handler.GetAllPermission)
	user.Get("/:id", handler.GetPermission)
}
