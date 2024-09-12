package utils

import (
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

// NotFoundUtils handles 404 - Route not found
func NotFoundUtil(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).JSON(struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{
		Success: false,
		Message: "Resource Not Found!",
	})
}

// ErrorUtil handles unhandled errors (500)
func ErrorUtil(c *fiber.Ctx) error {
	// Try to handle the request and capture any unhandled errors
	err := c.Next() // Process next middleware or route handler
	if err != nil {
		// Log the error for debugging
		log.Printf("Unhandled error: %v", err)

		// Return a 500 Internal Server Error response
		return c.Status(fiber.StatusNotFound).JSON(struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}{
			Success: false,
			Message: "Internal Server Error .. error",
		})
	}
	return nil
}

func RecoverWithLog() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic details
				log.Printf("Panic: %v\n", r)
				log.Printf("Stack trace: %s\n", string(debug.Stack()))

				// Send the panic error to ErrorHandler
				c.Status(fiber.StatusNotFound).JSON(struct {
					Success bool   `json:"success"`
					Message string `json:"message"`
				}{
					Success: false,
					Message: "Internal Server Error .. panic",
				})
			}
		}()

		return c.Next()
	}
}
