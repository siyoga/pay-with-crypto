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
	route.Get("/getForApprove", handlers.GetCardsForApprove) //TODO!: add middleware.AuthAdmin
	route.Post("/uploadLogo", middleware.Auth, handlers.CardLogoUploaderHandler)
	route.Post("/create", middleware.Auth, handlers.CardCreatorHandler)
	route.Post("/createTag", handlers.TagCreateHandler) //TODO!: add middleware.AuthAdmin
	route.Delete("/delete", middleware.Auth, handlers.CardDeleteHandler)
	route.Put("/edit", middleware.Auth, handlers.CardEditHandler)
	route.Get("/show/userid", handlers.CompanyGetByIdHandler)
	route.Get("/show", handlers.CardGetByIdHandler)
}
