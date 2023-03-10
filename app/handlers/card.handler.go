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
		return fiber.ErrInternalServerError
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
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func CardsSearcherByIdHandler(c *fiber.Ctx) error {
	var result []db.Card
	var state bool

	id := c.Query("id")

	if id == "" {
		return fiber.ErrBadRequest
	}

	if id != "" {
		result, state = db.SearchCardsById(id)
	}

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

func CardDeleteHandler(c *fiber.Ctx) error {
	var card db.Card
	var state bool
	loginedUser := c.Locals("user").(db.User).ID

	if err := c.BodyParser(&card); err != nil {
		return err
	}

	if !db.IsCardValidToLoginedUser(card.Id, loginedUser) {
		return fiber.ErrForbidden
	}

	if state = db.DeleteCardsById(card.Id); !state {
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

	if !db.IsCardValidToLoginedUser(changedCard.Id, loginedUser) {
		return fiber.ErrForbidden
	}

	db.UpdateCardOnId(changedCard)

	return c.Status(200).JSON(fiber.Map{"message": "Card successfully edited"})
}

func GetCardsForApprove(c *fiber.Ctx) error {
	cards, _ := db.GetManyBy[db.Card]("approved", false)

	return c.Status(200).JSON(cards)
}

func TagCreateHandler(c *fiber.Ctx) error {
	var newTag db.Tag
	admin := c.Locals("admin").(db.Admin)

	if err := c.BodyParser(&newTag); err != nil {
		return fiber.ErrBadRequest
	}

	newTag.ID = uuid.Must(uuid.NewV4())

	newTag.AdminID = admin.ID

	if ok := db.Add(newTag); !ok {
		return fiber.ErrInternalServerError
	}

	return c.Status(201).JSON(newTag)
}
