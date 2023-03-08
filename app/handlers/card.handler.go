package handlers

import (
	"os"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/datastore/s3"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

func CardSearcherHandler(c *fiber.Ctx) error {
	var card db.Card

	card.Name = c.Query("card_name")

	result, state := db.SearchCardByName(card.Name)

	if !state {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func CardLogoUploaderHandler(c *fiber.Ctx) error {
	logoBucket := os.Getenv("S3_BUCKET_CARD_LOGO")
	cardLogo, err := c.FormFile("cardLogo")
	cardId := c.Query("cardId")

	if cardId == "" {
		return fiber.ErrBadRequest
	}

	if err != nil {
		return fiber.ErrBadRequest
	}

	cardLogoBuffer, err := cardLogo.Open()

	if err != nil {
		return fiber.ErrBadRequest
	}

	defer cardLogoBuffer.Close()

	fileName := cardId + "_logo"
	fileNameInS3, isUploadOk := s3.UploadFile(cardLogo, cardLogoBuffer, logoBucket, fileName)

	if !isUploadOk {
		return fiber.ErrInternalServerError
	}

	_, isUpdateOk := db.UpdateOneBy[db.Card]("id", cardId, "image", *fileNameInS3)

	if !isUpdateOk {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func CardCreatorHandler(c *fiber.Ctx) error {
	var newCard db.Card
	user := c.Locals("user").(db.User)

	if err := c.BodyParser(&newCard); err != nil {
		return fiber.ErrBadRequest
	}

	newCard.Id = uuid.Must(uuid.NewV4()) // TODO: set userId from locals

	newCard.UserID = user.ID

	if ok := db.Add(newCard); !ok {
		return fiber.ErrInternalServerError
	}

	return c.Status(201).JSON(newCard)
}
