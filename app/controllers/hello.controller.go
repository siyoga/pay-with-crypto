package controllers

import (
	"pay-with-crypto/app/handlers"

	"github.com/gofiber/fiber/v2"
)

func HelloController(a *fiber.App) {
	route := a.Group("/hello")

	route.Get("/ping", handlers.PingHandler)
}