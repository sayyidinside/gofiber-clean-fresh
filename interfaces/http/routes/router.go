package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes/api"
)

func Setup(app *fiber.App) {
	cfg := config.AppConfig
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
		return fmt.Errorf("Test error")
	})

	apiGroupRoutes := app.Group("/api")
	api.SetupApiRoutes(apiGroupRoutes)
}
