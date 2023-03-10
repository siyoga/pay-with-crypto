package handlers

import (
	db "pay-with-crypto/app/datastore"

	"github.com/gofiber/fiber/v2"
)

func CompanyShowByIdHandler(c *fiber.Ctx) error {
	var company db.User
	var state bool

	if err := c.BodyParser(&company); err != nil {
		return fiber.ErrBadRequest
	}

	company, state = db.GetOneBy[db.User]("id", company.ID)

	if !state {
		return fiber.ErrNotFound
	}

	return c.Status(fiber.StatusOK).JSON(company)
}
