package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
)

func RegisterRoutes(route fiber.Router, handler handler.AuthHandler) {
	authRoutes := route.Group("/auth")

	authRoutes.Post("/login", handler.Login)
}
