package handlers

import (
	"os"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/datastore/s3"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

func CardSearcherByNameHandler(c *fiber.Ctx) error {
	var result []db.Card
	var state bool

	value := c.Query("card_name")

	if value != "" {
		result, state = db.SearchCardByName(value)
	}

	if !state {
		return fiber.ErrNotFound
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func CardsSearcherByTagsHandler(c *fiber.Ctx) error {
	var result []db.Card
	var state bool

	value := c.Query("tags")

	if value != "" {
		result, state = db.SearchCardsByTags(value)
	}

	if !state {
		return fiber.ErrNotFound
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func CardsSearcherByIdHandler(c *fiber.Ctx) error {
	var result db.Card
	var state bool

	id := c.Query("id")

	if id == "" {
		return fiber.ErrBadRequest
	}

	if id != "" {
		result, state = db.GetOneBy[db.Card]("id", id)
	}

	if !state {
		return fiber.ErrNotFound
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

	newCard.ID = uuid.Must(uuid.NewV4())

	newCard.UserID = user.ID

	if ok := db.Add(newCard); !ok {
		return fiber.ErrInternalServerError
	}

	return c.Status(201).JSON(newCard)
}

func CardDeleteHandler(c *fiber.Ctx) error {
	var card db.Card
	var state bool
	loginedUser := c.Locals("user").(db.User).ID

	if err := c.BodyParser(&card); err != nil {
		return err
	}

	if !db.IsCardValidToLoginedUser(card.ID, loginedUser) {
		return fiber.ErrForbidden
	}

	if state = db.DeleteBy[db.Card]("id", card.ID); !state {
		return fiber.ErrNotFound
	}

	return c.Status(fiber.StatusOK).JSON(state)
}

func CardEditHandler(c *fiber.Ctx) error {
	var changedCard db.Card
	loginedUser := c.Locals("user").(db.User).ID

	if err := c.BodyParser(&changedCard); err != nil {
		return fiber.ErrBadRequest
	}

	if _, found := db.GetOneBy[db.Card]("id", changedCard.ID); !found {
		return fiber.ErrBadRequest
	}

	if !db.IsCardValidToLoginedUser(changedCard.ID, loginedUser) {
		return fiber.ErrForbidden
	}

	if !db.WholeOneUpdate(changedCard) {
		return fiber.ErrInternalServerError
	}

	return c.Status(200).JSON(fiber.Map{"message": "Card successfully edited"})
}

func CardGetByIdHandler(c *fiber.Ctx) error {
	var card db.Card
	var state bool

	cardId := c.Query("id")

	if cardId == "" {
		return fiber.ErrBadRequest
	}

	card, state = db.GetOneBy[db.Card]("id", cardId)

	if !state {
		return fiber.ErrNotFound
	}

	return c.Status(fiber.StatusOK).JSON(card)
}
