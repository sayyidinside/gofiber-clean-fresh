package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
)

func RegisterModuleRoutes(route fiber.Router, handler handler.ModuleHandler) {
	user := route.Group("/modules")

	user.Get("/", handler.GetAllModule)
	user.Get("/:id", handler.GetModule)
	user.Post("", handler.CreateModule)
	user.Put("/:id", handler.UpdateModule)
	user.Delete("/:id", handler.DeleteModule)
}
