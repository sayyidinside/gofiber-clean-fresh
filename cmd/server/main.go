package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/cmd/bootstrap"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/config"
	"github.com/sayyidinside/gofiber-clean-fresh/infrastructure/database"
	"github.com/sayyidinside/gofiber-clean-fresh/pkg/utils"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Println(err.Error())
	}

	app := fiber.New()

	// Recover panic
	app.Use(utils.RecoverWithLog())

	app.Use(utils.ErrorUtil)

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	bootstrap.Initialize(app, db)

	app.Use(utils.NotFoundUtil)

	app.Listen(fmt.Sprintf(":%s", config.AppConfig.Port))
}
