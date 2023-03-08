package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func CardController(a *fiber.App) {
	route := a.Group("/card")

	route.Get("/search", middleware.Auth, handlers.CardSearcherHandler)
	route.Post("/uploadLogo", middleware.Auth, handlers.CardLogoUploaderHandler)
	route.Post("/create", middleware.Auth, handlers.CardCreatorHandler)
	route.Put("/edit", handlers.CardEditHandler)
}
