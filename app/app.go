package app

import (
	"fmt"
	"pay-with-crypto/app/controllers"
	db "pay-with-crypto/app/datastore"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Start(config db.DatabaseConfig) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-api-key, accessToken",
		AllowCredentials: true,
	}))

	controllers.AuthController(app)
	controllers.CardController(app)
	controllers.CompanyController(app)
	controllers.TagController(app)

	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.User, config.Password, config.Database))

	return app
}
