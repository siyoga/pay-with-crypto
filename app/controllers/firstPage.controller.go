package controllers

import (
	"pay-with-crypto/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func FirstPageController(a *fiber.App) {
	route := a.Group("/")

	route.Get("/", handlers.ShowFirstCards)
}
