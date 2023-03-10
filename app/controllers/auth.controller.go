package controllers

import (
	"pay-with-crypto/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func AuthController(a *fiber.App) {
	route := a.Group("/auth")

	route.Post("/register", handlers.RegisterHandler)
	route.Post("/login", handlers.LoginHandler)
	route.Post("/admin_register", handlers.AdminRegisterHandler)
}
