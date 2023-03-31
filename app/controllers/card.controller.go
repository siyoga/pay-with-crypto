package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func CardController(a *fiber.App) {
	route := a.Group("/card")

	a.Get("/", handlers.ShowFirstCards)
	route.Get("/search", handlers.CardsSearcher)
	route.Get("/search/id", handlers.CardsSearcherByIdHandler)
	route.Post("/uploadLogo", middleware.Auth, handlers.CardLogoUploaderHandler)
	route.Get("/getLogo", handlers.CardLogoGetterHandler)
	route.Post("/create", middleware.Auth, handlers.CardCreatorHandler)
	route.Delete("/delete", middleware.Auth, handlers.CardDeleteHandler)
	route.Put("/edit", middleware.Auth, handlers.CardEditHandler)
}
