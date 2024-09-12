package v1

import "github.com/gofiber/fiber/v2"

func SetupAPIV1Routes(v1 fiber.Router) {
	v1.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: true,
			Message: "API is running smoothly",
		})
	})
}
