package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthController(a *fiber.App) {
	route := a.Group("/auth")

	route.Get("/google/register", handlers.AuthGoogleGetApprove)
	route.Get("/google/callback", handlers.Callback)
	route.Get("/google/login", handlers.AuthGoogleGetApprove) //TODO!: Separate login from register. Now when register we login
	route.Get("/token_update", middleware.Auth, handlers.GetValidTokensHandler)
	route.Post("/register", handlers.RegisterHandler)
	route.Post("/login", handlers.LoginHandler)
	route.Post("/admin_register", middleware.AuthAdmin, handlers.AdminRegisterHandler)
	route.Post("/admin_login", handlers.AdminLoginHandler)
}
