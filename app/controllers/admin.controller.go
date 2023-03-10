package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func AdminController(a *fiber.App) {
	route := a.Group("/admin")

	route.Get("/getForApprove", middleware.AuthAdmin, handlers.GetCardsForApprove)
	route.Post("/createTag", middleware.AuthAdmin, handlers.TagCreateHandler)
}
