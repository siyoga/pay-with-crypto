package handlers

import (
	db "pay-with-crypto/app/datastore"

	"github.com/gofiber/fiber/v2"
)

func CardSearcherHandler(c *fiber.Ctx) error {
	var card db.Card

	card.Name = c.Query("card_name")

	if err := c.BodyParser(&card); err != nil {
		return fiber.ErrBadRequest
	}

	result, state := db.SearchCardByName(card.Name)

	if !state {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
