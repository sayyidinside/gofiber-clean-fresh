package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handlers"
	v1 "github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes/v1"
)

func Setup(app *fiber.App, handler *handlers.Handlers) {
	api := app.Group("/api")
	v1.RegisterRoutes(api, handler)

	cfg := config.AppConfig
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("App is running, with name %s", cfg.AppName))
	})

	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("This panic is caught by fiber")
	})
}
