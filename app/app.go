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
	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.User, config.Password, config.Database))

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,X-Requested-With,Content-Type,Accept,Authorization",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	controllers.AuthController(app)
	controllers.CardController(app)
	controllers.CompanyController(app)
	controllers.TagController(app)

	return app
}
