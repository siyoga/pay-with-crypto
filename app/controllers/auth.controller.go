package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthController(a *fiber.App) {
	route := a.Group("/auth")

	route.Get("/whoami", handlers.WhoAmIHandler)
	route.Post("/tokenUpdate", handlers.UpdateTokensHandler)
	route.Post("/register", handlers.RegisterHandler)
	route.Post("/login", handlers.LoginHandler)
	route.Post("/admin/register", middleware.AuthAdmin, handlers.AdminRegisterHandler)
	route.Post("/admin/login", handlers.AdminLoginHandler)
	route.Post("/google/register", handlers.GoogleRegisterHandler)
}
