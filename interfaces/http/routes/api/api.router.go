package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes/api/tests"
	v1 "github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes/api/v1"
)

func SetupApiRoutes(api fiber.Router) {
	api.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: true,
			Message: "API is running smoothly",
		})
	})

	apiV1GroupRoutes := api.Group("/v1")
	v1.SetupAPIV1Routes(apiV1GroupRoutes)

	apiTestGroupRoutes := api.Group("/tests")
	tests.SetupApiTestRoutes(apiTestGroupRoutes)
}
