package controllers

import (
	"pay-with-crypto/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func AuthController(a *fiber.App) {
	route := a.Group("/auth")

	route.Get("/whoami", handlers.WhoAmIHandler)
	route.Post("/refresh", handlers.RefreshHandler)
	route.Post("/register", handlers.RegisterHandler)
	route.Post("/login", handlers.LoginHandler)
}
