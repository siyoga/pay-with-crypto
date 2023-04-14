package controllers

import (
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func AdminController(a *fiber.App) {
	route := a.Group("/admin")

	route.Get("/getForApprove", middleware.AuthAdmin, handlers.GetCardsForApprove)
	route.Get("/getAllCompanies", middleware.AuthAdmin, handlers.GetAllCompanies)
	route.Put("/validateCard", middleware.AuthAdmin, handlers.ValidateCard)
	route.Put("/validateTag", middleware.AuthAdmin, handlers.ValidateTag)
	route.Patch("/unban", middleware.AuthAdmin, handlers.UnbanCompanyHandler) //TODO!: Need to do a personal list of banned for every admin as a source of banned id of Companys
	route.Delete("/ban", middleware.AuthAdmin, handlers.SoftDeleteHandler)    //
	route.Delete("/deleteTag", middleware.AuthAdmin, handlers.TagDeleteHandler)
}
