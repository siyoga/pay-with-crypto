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
	route.Get("/whoami", handlers.WhoAmIHandler)
	route.Get("/google/login", handlers.AuthGoogleGetApprove) //TODO!: Separate login from register. Now when register we login
	route.Post("/tokenUpdate", handlers.UpdateTokensHandler)
	route.Post("/register", handlers.RegisterHandler)
	route.Post("/login", handlers.LoginHandler)
	route.Post("/admin/register", middleware.AuthAdmin, handlers.AdminRegisterHandler)
	route.Post("/admin/login", handlers.AdminLoginHandler)
	route.Post("/google/tokeninfo", handlers.TokenInfo)
}
