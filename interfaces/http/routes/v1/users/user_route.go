package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handler"
)

func RegisterUserRoutes(route fiber.Router, handler handler.UserHandler) {
	user := route.Group("/informations")

	user.Get("/:id", handler.GetUser)
	// user.Get("/:id", func(c *fiber.Ctx) error {
	// 	// Extract identifier and username
	// 	ctx := helpers.ExtractIdentifierAndUsername(c)

	// 	// Create initial log system
	// 	logData := helpers.CreateLogSystem(ctx, "Get User by ID Request")

	// 	// Call the handler with the new context
	// 	return handler.GetUser(c, ctx, logData)
	// })

	user.Get("/", handler.GetAllUser)
	user.Post("/", handler.CreateUser)

	user.Put("/:id/reset-password", handler.ResetPassword)
	user.Put("/:id", handler.UpdateUser)

	user.Delete("/:id", handler.DeleteUser)
}
