package handlers

import (
	db "pay-with-crypto/app/datastore"

	"github.com/gofiber/fiber/v2"
)

func CompanyGetByIdHandler(c *fiber.Ctx) error {
	var company db.User
	var state bool

	companyId := c.Query("id")

	if companyId == "" {
		return fiber.ErrBadRequest
	}

	company, state = db.GetOneBy[db.User]("id", companyId)

	if !state {
		return fiber.ErrNotFound
	}

	return c.Status(fiber.StatusOK).JSON(company)
}
