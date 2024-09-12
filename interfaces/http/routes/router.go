package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes/api"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

func Setup(app *fiber.App) {
	cfg := config.AppConfig
	app.Get("/", func(c *fiber.Ctx) error {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  200,
			Success: true,
			Message: fmt.Sprintf("App is running, with name %s", cfg.AppName),
		})
	})

	apiGroupRoutes := app.Group("/api")
	api.SetupApiRoutes(apiGroupRoutes)
}
