package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handlers"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/middleware"
	v1 "github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes/v1"
)

func Setup(app *fiber.App, handler *handlers.Handlers) {
	cfg := config.AppConfig

	api := app.Group("/api")

	// Apply middleware for general API routes
	api.Use(middleware.CORS())

	v1.RegisterRoutes(api, handler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: true,
			Message: fmt.Sprintf("App is running, with name %s", cfg.AppName),
		})
	})

	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("This panic is caught by fiber")
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		return fmt.Errorf("test error")
	})
}
