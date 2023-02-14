package app

import (
	"pay-with-crypto/app/controllers"
	db "pay-with-crypto/app/datastore"

	"github.com/gofiber/fiber/v2"
)

func Start(config db.DatabaseConfig) *fiber.App {
	app := fiber.New()

	controllers.AuthController(app)

	return app
}
