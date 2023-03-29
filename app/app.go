package app

import (
	"fmt"
	"pay-with-crypto/app/controllers"
	db "pay-with-crypto/app/datastore"

	"github.com/gofiber/fiber/v2"
)

func Start(config db.DatabaseConfig) *fiber.App {
	app := fiber.New()

	controllers.AuthController(app)
	controllers.CardController(app)
	controllers.CompanyController(app)
	controllers.AdminController(app)

	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.User, config.Password, config.Database))

	return app
}
