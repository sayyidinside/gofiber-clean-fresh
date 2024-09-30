package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
)

func RegisterUserRoutes(route fiber.Router, handler handler.UserHandler) {
	user := route.Group("/informations")

	user.Get("/:id", handler.GetUser)
	user.Get("/", handler.GetAllUser)
	user.Post("/", handler.CreateUser)
}
