package handlers

import (
	db "pay-with-crypto/app/datastore"

	"github.com/gofiber/fiber/v2"
)

func CardSearcherHandler(c *fiber.Ctx) error {
	var card db.Card

	if err := c.BodyParser(&card); err != nil {
		return fiber.ErrBadRequest
	}

	resultOfSearch, statementOfSearch := db.SearchCardByName(card.Name)

	if !statementOfSearch {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(resultOfSearch)
}
