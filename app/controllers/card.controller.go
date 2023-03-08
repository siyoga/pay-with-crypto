package controllers

import (
	"pay-with-crypto/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func CardController(a *fiber.App) {
	route := a.Group("/card")

	route.Get("/search", handlers.CardSearcherHandler)

	route.Put("/edit", handlers.CardEditHandler)
}
