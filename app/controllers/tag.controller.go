package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func TagController(a *fiber.App) {
	route := a.Group("/tag")

	route.Get("/get/all", handlers.TagGetterHandler)
	route.Post("/create", middleware.Auth, handlers.TagCreateHandler)
}
