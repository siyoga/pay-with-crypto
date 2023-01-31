package app

import (
	"log"
	"pay-with-crypto/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Start(port string) {
	app := fiber.New()

	controllers.HelloController(app)

	if err:= app.Listen(":" + port); err != nil {
		log.Panic(err)
	}
}