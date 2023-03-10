package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func CompanyController(a *fiber.App) {
	route := a.Group("/company")

	route.Get("/search/userid", handlers.CompanyGetByIdHandler)
	route.Post("/uploadLogo", middleware.Auth, handlers.CompanyLogoUploaderHandler)
}
