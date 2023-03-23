package handlers

import (
	"fmt"
	"os"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/datastore/s3"

	"github.com/go-ping/ping"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

func CardsSearcher(c *fiber.Ctx) error {
	var result []db.Card
	var state bool

	name := c.Query("name")
	tags := c.Query("tags")

	if name != "" && tags != "" {
		result, state = db.SearchCard(name, tags)
	} else if name != "" {
		result, state = db.SearchCardByName(name)
	} else if tags != "" {
		result, state = db.SearchCardsByTags(tags)
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
		if result, state = db.GetOneBy[db.Card]("id", id); !state {
			return fiber.ErrNotFound
		}
		result.Views++
		if !db.WholeOneUpdate(result) {
			return fiber.ErrInternalServerError
		}
	}

	if db.IsCardOwnerSoftDeleted(result.CompanyID) {
		return fiber.ErrForbidden
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func CardLogoGetterHandler(c *fiber.Ctx) error {
	cardId := c.Query("cardId")

	if cardId == "" {
		return fiber.ErrBadRequest
	}

	card, state := db.GetOneBy[db.Card]("id", cardId)

	if !state {
		return fiber.ErrNotFound
	}

	var output string

	if card.Image == "" {
		output = ""
	} else {
		output = fmt.Sprintf("http://217.25.95.4:9000/card-logos/%s", card.Image)
	}

	return c.JSON(fiber.Map{"link": output})
}

func CardLogoUploaderHandler(c *fiber.Ctx) error {
	logoBucket := os.Getenv("S3_BUCKET_CARD_LOGO")
	cardLogo, err := c.FormFile("cardLogo")
	cardId := c.Query("cardId")
	var card db.Card
	var state bool

	if cardId == "" {
		return fiber.ErrBadRequest
	}

	if card, state = db.GetOneBy[db.Card]("id", cardId); !state {
		return fiber.ErrNotFound
	}

	if db.IsCardOwnerSoftDeleted(card.CompanyID) {
		return fiber.ErrForbidden
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
	company := c.Locals("company").(db.Company)

	if err := c.BodyParser(&newCard); err != nil {
		return fiber.ErrBadRequest
	}

	if _, engaged := db.GetOneBy[db.Card]("name", newCard.Name); engaged {
		return fiber.ErrConflict
	}

	pinger, err := ping.NewPinger(newCard.LinkToProd)
	if err != nil {
		return fiber.ErrBadRequest
	}
	pinger.Count = 3
	pinger.TTL = 129
	err = pinger.Run() // Blocks until finished

	newCard.ID = uuid.Must(uuid.NewV4())
	newCard.CompanyID = company.ID
	newCard.Approved = "pending"

	if ok := db.Add(newCard); !ok {
		return fiber.ErrInternalServerError
	}

	return c.Status(201).JSON(newCard)
}

func CardDeleteHandler(c *fiber.Ctx) error {
	var card db.Card
	var state bool
	loginedUser := c.Locals("company").(db.Company).ID

	if err := c.BodyParser(&card); err != nil {
		return err
	}

	if card, state = db.GetOneBy[db.Card]("id", card.ID); !state {
		return fiber.ErrBadRequest
	}

	if db.IsCardOwnerSoftDeleted(card.CompanyID) {
		return fiber.ErrForbidden
	}

	if !db.IsValid(card.CompanyID, loginedUser) {
		return fiber.ErrForbidden
	}

	if state = db.DeleteBy[db.Card]("id", card.ID); !state {
		return fiber.ErrNotFound
	}

	return c.Status(fiber.StatusOK).JSON(state)
}

func CardEditHandler(c *fiber.Ctx) error {
	var changedCard db.Card
	var state bool
	loginedCompany := c.Locals("company").(db.Company).ID

	if err := c.BodyParser(&changedCard); err != nil {
		return fiber.ErrBadRequest
	}

	if changedCard, state = db.GetOneBy[db.Card]("id", changedCard.ID); !state {
		return fiber.ErrBadRequest
	}

	if db.IsCardOwnerSoftDeleted(changedCard.CompanyID) {
		return fiber.ErrForbidden
	}

	if !db.IsValid(changedCard.CompanyID, loginedCompany) {

		return c.Status(200).JSON(changedCard.CompanyID)
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

	if card, state = db.GetOneBy[db.Card]("id", cardId); !state {
		return fiber.ErrNotFound
	}

	if db.IsCardOwnerSoftDeleted(card.CompanyID) {
		return fiber.ErrForbidden
	}

	return c.Status(fiber.StatusOK).JSON(card)
}
