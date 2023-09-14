package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func CardController(a *fiber.App) {
	route := a.Group("/card")

	route.Get("/get", handlers.GetCardsByCompany)
	route.Get("/get/all", handlers.CardGetAll)
	route.Post("/logo/upload", middleware.Auth, handlers.CardLogoUploadHandler)
	route.Post("/create", middleware.Auth, handlers.CardCreatorHandler)
	route.Delete("/delete", middleware.Auth, handlers.CardDeleteHandler)
	route.Put("/edit", middleware.Auth, handlers.CardEditHandler)
}
