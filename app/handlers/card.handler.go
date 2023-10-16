package handlers

import (
	"fmt"
	"path/filepath"
	d "pay-with-crypto/app/datastore"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/datastore/s3"
	"pay-with-crypto/app/utility"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

func CardLogoUploadHandler(c *fiber.Ctx) error {
	cardLogoRaw, err := c.FormFile("cardLogo")
	cardId := c.Get("Card")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request"})
	}

	card, isOk := d.GetOneBy[d.Card]("id", cardId)

	if !isOk {
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Card not found"})
	}

	cardLogo, err := cardLogoRaw.Open()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request, provide companyLogo correctly"})
	}

	defer cardLogo.Close()

	fileName := card.ID.String() + "_logo" + filepath.Ext(cardLogoRaw.Filename)
	s3.DeleteImage(fileName) // чтобы обезопаситься от дубликатов

	fileLink, isUploadOk := s3.UploadFile(cardLogo, fileName)

	if !isUploadOk {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	card, isUpdateOk := db.UpdateOneBy[db.Card]("id", card.ID, "logoLink", fileLink)

	if !isUpdateOk {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.Status(fiber.StatusOK).JSON(card)
}

func CardCreatorHandler(c *fiber.Ctx) error {
	var newCard db.Card
	company := c.Locals("company").(db.Company)

	if err := c.BodyParser(&newCard); err != nil {
		fmt.Println("err")
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Неверный запрос"})
	}

	if _, exist := db.GetOneBy[db.Card]("name", newCard.Name); exist {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "Карточка с таким именем уже существует"})
	}

	newCard.ID = uuid.Must(uuid.NewV4())
	newCard.CompanyOwner = company.ID
	newCard.Approved = "pending"

	if ok := db.Add(newCard); !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Что-то пошло не так, попробуйте позже"})
	}

	return c.Status(201).JSON(newCard)
}

func GetCardsByCompany(c *fiber.Ctx) error {
	companyId := c.Query("id")

	cards, isOk := db.GetManyBy[db.Card]("company_owner", companyId)

	if !isOk {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Что-то пошло не так, попробуйте позже"})
	}

	return c.Status(fiber.StatusOK).JSON(cards)
}

func CardGetAll(c *fiber.Ctx) error {
	cards, isOk := db.GetManyBy[db.Card]("approved", "pending")

	if !isOk {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Что-то пошло не так, попробуйте позже"})
	}

	return c.Status(fiber.StatusOK).JSON(cards)
}

func CardDeleteHandler(c *fiber.Ctx) error {
	var card db.Card
	var state bool
	loginedUser := c.Locals("company").(db.Company).ID

	if err := c.BodyParser(&card); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request"})
	}

	if card, state = db.GetOneBy[db.Card]("id", card.ID); !state {
		if _, state = db.GetOneUnscopedBy[db.Card]("id", card.ID); state {
			return c.Status(fiber.StatusForbidden).JSON(utility.Message{Text: "Owner of card was banned"})
		}
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Card not exist"})
	}

	if !db.IsValid(card.CompanyOwner, loginedUser) {
		return c.Status(fiber.StatusForbidden).JSON(utility.Message{Text: "This card not belongs to your company"})
	}

	if state = db.DeleteBy[db.Card]("id", card.ID); !state {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.SendStatus(204)
}

func CardEditHandler(c *fiber.Ctx) error {
	var changedCard db.Card
	var state bool
	loginedCompany := c.Locals("company").(db.Company).ID

	if err := c.BodyParser(&changedCard); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request"})
	}

	if changedCard, state = db.GetOneBy[db.Card]("id", changedCard.ID); !state {
		if _, state = db.GetOneUnscopedBy[db.Card]("id", changedCard.ID); state {
			return c.Status(fiber.StatusForbidden).JSON(utility.Message{Text: "Owner of card was banned"})
		}
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Card not exist"})
	}

	if !db.IsValid(changedCard.CompanyOwner, loginedCompany) {
		return c.Status(fiber.StatusForbidden).JSON(utility.Message{Text: "This card not belongs to your company"})
	}

	if !db.WholeOneUpdate(changedCard) {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.Status(200).JSON(utility.Message{Text: "Card successfully edited"})
}
