package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func CompanyController(a *fiber.App) {
	route := a.Group("/company")

	route.Get("/search/id", handlers.GetByIdHandler)
	route.Post("/update", middleware.Auth, handlers.UpdateHandler)
	route.Post("/uploadLogo", middleware.Auth, handlers.LogoUploadHandler)
	route.Post("/createTag", middleware.Auth, handlers.TagCreateHandler)
}
