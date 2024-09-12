package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handlers"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes/tests"
	v1 "github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes/v1"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

func Setup(app *fiber.App, handler *handlers.Handlers) {
	api := app.Group("/api")
	v1.RegisterRoutes(api, handler)

	test := app.Group("/tests")
	tests.SetupApiTestRoutes(test)

	cfg := config.AppConfig
	app.Get("/", func(c *fiber.Ctx) error {
		return helpers.ResponseFormatter(c, helpers.BaseResponse{
			Status:  200,
			Success: true,
			Message: fmt.Sprintf("App is running, with name %s", cfg.AppName),
		})
	})
}
