package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/database"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Println(err.Error())
	}

	app := fiber.New()

	// Recover panic
	app.Use(recover.New())

	_, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	routes.Setup(app)

	app.Listen(fmt.Sprintf(":%s", config.AppConfig.Port))
}
