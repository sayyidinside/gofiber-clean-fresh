package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
)

func Setup(app *fiber.App) {
	cfg := config.AppConfig
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("App is running, with name %s", cfg.AppName))
	})

	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("This panic is caught by fiber")
	})
}
