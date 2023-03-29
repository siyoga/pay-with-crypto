package handlers

import (
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/utility"

	"github.com/gofiber/fiber/v2"
)

func ShowFirstCards(c *fiber.Ctx) error {
	result, state := db.GetAllOrdered[db.Card]("approved", "approved", "views desc")
	if !state {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "No cards!"})
	}

	return (c.Status(fiber.StatusOK).JSON(result))
}
