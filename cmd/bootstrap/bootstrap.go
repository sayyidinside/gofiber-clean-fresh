package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayyidinside/gofiber-clean-fresh/domain/user"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/handlers"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/routes"
	userHand "github.com/sayyidinside/gofiber-clean-fresh/interfaces/http/user"
	"gorm.io/gorm"
)

func Initialize(app *fiber.App, db *gorm.DB) {
	// Repositories
	userRepo := user.NewRepository(db)

	// Service
	userService := user.NewService(userRepo)

	// Handler
	userHandler := userHand.NewHandler(userService)

	// Setup handler to send to routes setup
	handler := &handlers.Handlers{
		UserManagementHandler: &handlers.UserManagementHandler{
			UserHandler: *userHandler,
		},
	}

	routes.Setup(app, handler)
}
