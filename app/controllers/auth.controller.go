package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthController(a *fiber.App) {
	route := a.Group("/auth")

	route.Post("/register", middleware.Auth, handlers.RegisterHandler)
	route.Post("/login", middleware.Auth, handlers.LoginHandler)
}
