package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/cmd/bootstrap"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/database"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/helpers"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Println(err.Error())
	}

	app := fiber.New()

	// Recover panic
	app.Use(helpers.RecoverWithLog())

	app.Use(helpers.ErrorHelper)

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	bootstrap.Initialize(app, db)

	app.Use(helpers.NotFoundHelper)

	app.Listen(fmt.Sprintf(":%s", config.AppConfig.Port))
}
