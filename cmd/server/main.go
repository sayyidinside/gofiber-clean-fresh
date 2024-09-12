package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sayyidinside/gofiber-clean-fresh/cmd/bootstrap"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/database"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Println(err.Error())
	}

	app := fiber.New()

	// Recover panic
	app.Use(recover.New())

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	bootstrap.Initialize(app, db)

	app.Listen(fmt.Sprintf(":%s", config.AppConfig.Port))
}
