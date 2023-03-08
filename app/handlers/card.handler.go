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

func CardUpdateHandler(c *fiber.Ctx) error {
	var changedCard db.Card

	if err := c.BodyParser(&changedCard); err != nil {
		return fiber.ErrBadRequest
	}

	result, ok := db.GetOneBy[db.Card]("id", changedCard.Id)
	if ok == false {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Not found card"})
	}

	db.Datastore.Model(&result).Updates(changedCard)

	return c.Status(200).JSON(fiber.Map{"message": "Card successfully edited"})
}
