package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func CardController(a *fiber.App) {
	route := a.Group("/card")

	route.Get("/search", handlers.CardSearcherByNameHandler)
	route.Get("/search/tags", handlers.CardsSearcherByTagsHandler)
	route.Get("/search/id", handlers.CardsSearcherByIdHandler)
	route.Post("/uploadLogo", middleware.Auth, handlers.CardLogoUploaderHandler)
	route.Post("/create", middleware.Auth, handlers.CardCreatorHandler)
	route.Put("/edit", middleware.Auth, handlers.CardEditHandler)

}
