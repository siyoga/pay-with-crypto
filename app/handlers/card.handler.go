package handlers

import (
	"fmt"
	"os"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/datastore/s3"
	"pay-with-crypto/app/utility"

	"github.com/go-ping/ping"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

// @Description Search card
// @Tags Card
// @Accept json
// @Produce json
// @Param name query string false "Card name"
// @Param tags query string false "Card tags"
// @Success 200 {object} []datastore.Card
// @Failure 404 {object} utility.Message "No cards"
// @Router /card/search [get]
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
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "No cards with appropriate parameters found"})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// @Description Search card by id
// @Tags Card
// @Accept json
// @Produce json
// @Param cardId query string true "Card id"
// @Success 200 {object} datastore.Card
// @Failure 400 {object} utility.Message "Invalid request"
// @Failure 404 {object} utility.Message "No card"
// @Failure 403 {object} utility.Message "Card owner was banned"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /card/search/id [get]
func CardsSearcherByIdHandler(c *fiber.Ctx) error {
	var result db.Card
	var state bool

	id := c.Query("cardId")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request"})
	}

	if result, state = db.GetOneBy[db.Card]("id", id); !state {
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Card not exist"})
	}

	result.Views++
	if !db.WholeOneUpdate(result) {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	if db.IsCardOwnerSoftDeleted(result.CompanyID) {
		return c.Status(fiber.StatusForbidden).JSON(utility.Message{Text: "The account holding that card has been deleted."})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// @Description Get card logo
// @Tags Card
// @Accept json
// @Produce json
// @Param cardId query string true "Card id"
// @Success 200 {object} utility.Message "Card logo link"
// @Failure 400 {object} utility.Message "Invalid request"
// @Failure 404 {object} utility.Message "No card"
// @Router /card/getLogo [get]
func CardLogoGetterHandler(c *fiber.Ctx) error {
	cardId := c.Query("cardId")

	if cardId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request"})
	}

	card, state := db.GetOneBy[db.Card]("id", cardId)

	if !state {
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Card not exist"})
	}

	var output string

	if card.Image == "" {
		output = ""
	} else {
		output = fmt.Sprintf("http://217.25.95.4:9000/card-logos/%s", card.Image)
	}

	return c.Status(200).JSON(utility.Message{Text: output})
}

// @Description Card logo uploader
// @Tags Card
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param cardId query string true "Card id"
// @Param cardLogo formData file true "Logo image"
// @Success 204 "Card successful uploaded"
// @Failure 400 {object} utility.Message "Invalid request"
// @Failure 400 {object} utility.Message "Invalid image"
// @Failure 403 {object} utility.Message "Card owner was banned"
// @Failure 404 {object} utility.Message "No card"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /card/uploadLogo [post]
func CardLogoUploaderHandler(c *fiber.Ctx) error {
	logoBucket := os.Getenv("S3_BUCKET_CARD_LOGO")
	cardLogo, err := c.FormFile("cardLogo")
	cardId := c.Query("cardId")
	var card db.Card
	var state bool

	if cardId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request, provide cardId"})
	}

	if card, state = db.GetOneBy[db.Card]("id", cardId); !state {
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Card not exist"})
	}

	if db.IsCardOwnerSoftDeleted(card.CompanyID) {
		return c.Status(fiber.StatusForbidden).JSON(utility.Message{Text: "The account holding that card has been deleted."})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request, provide cardLogo"})
	}

	cardLogoBuffer, err := cardLogo.Open()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid image"})
	}

	defer cardLogoBuffer.Close()

	fileName := cardId + "_logo"
	fileNameInS3, isUploadOk := s3.UploadFile(cardLogo, cardLogoBuffer, logoBucket, fileName)

	if !isUploadOk {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	_, isUpdateOk := db.UpdateOneBy[db.Card]("id", cardId, "image", *fileNameInS3)

	if !isUpdateOk {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// @Description Card create
// @Tags Card
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param card_data body object{name=string,description=string,price=string,linkToProd=string,tags=[]string} true "Card data"
// @Success 201 {object} datastore.Card "Card successful created"
// @Failure 400 {object} utility.Message "Invalid request"
// @Failure 409 {object} utility.Message "Already created"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /card/create [post]
func CardCreatorHandler(c *fiber.Ctx) error {
	var newCard db.Card
	company := c.Locals("company").(db.Company)

	if err := c.BodyParser(&newCard); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request"})
	}

	if _, engaged := db.GetOneBy[db.Card]("name", newCard.Name); engaged {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "Card with that name already exist"})
	}

	pinger, err := ping.NewPinger(newCard.LinkToProd)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid link to company."})
	}
	pinger.Count = 3
	pinger.TTL = 129
	err = pinger.Run() // Blocks until finished
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Transmitted link does not respond."})
	}

	newCard.ID = uuid.Must(uuid.NewV4())
	newCard.CompanyID = company.ID
	newCard.Approved = "pending"

	if ok := db.Add(newCard); !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.Status(201).JSON(newCard)
}

// @Description Card delete
// @Tags Card
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param card_data body object{id=string} true "Card data"
// @Success 204 "Card successful deleted"
// @Failure 400 {object} utility.Message "Invalid request"
// @Failure 403 {object} utility.Message "Card owner was banned"
// @Failure 403 {object} utility.Message "Other owner"
// @Failure 404 {object} utility.Message "No card"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /card/delete [delete]
func CardDeleteHandler(c *fiber.Ctx) error {
	var card db.Card
	var state bool
	loginedUser := c.Locals("company").(db.Company).ID

	if err := c.BodyParser(&card); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request"})
	}

	if card, state = db.GetOneBy[db.Card]("id", card.ID); !state {
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Card not exist"})
	}

	if db.IsCardOwnerSoftDeleted(card.CompanyID) {
		return c.Status(fiber.StatusForbidden).JSON(utility.Message{Text: "The account holding that card has been deleted."})
	}

	if !db.IsValid(card.CompanyID, loginedUser) {
		return c.Status(fiber.StatusForbidden).JSON(utility.Message{Text: "This card not belongs to your company"})
	}

	if state = db.DeleteBy[db.Card]("id", card.ID); !state {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.SendStatus(204)
}

// @Description Card edit
// @Tags Card
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param card_data body object{id=string,name=string,description=string,price=string,linkToProd=string,tags=[]string} true "Card data"
// @Success 200 {object} utility.Message "Card successful edited"
// @Failure 400 {object} utility.Message "Invalid request"
// @Failure 403 {object} utility.Message "Card owner was banned"
// @Failure 403 {object} utility.Message "Other owner"
// @Failure 404 {object} utility.Message "No card"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /card/edit [put]
func CardEditHandler(c *fiber.Ctx) error {
	var changedCard db.Card
	var state bool
	loginedCompany := c.Locals("company").(db.Company).ID

	if err := c.BodyParser(&changedCard); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request"})
	}

	if changedCard, state = db.GetOneBy[db.Card]("id", changedCard.ID); !state {
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Card not exist"})
	}

	if db.IsCardOwnerSoftDeleted(changedCard.CompanyID) {
		return c.Status(fiber.StatusForbidden).JSON(utility.Message{Text: "The account holding that card has been deleted."})
	}

	if !db.IsValid(changedCard.CompanyID, loginedCompany) {
		return c.Status(fiber.StatusForbidden).JSON(utility.Message{Text: "This card not belongs to your company"})
	}

	if !db.WholeOneUpdate(changedCard) {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.Status(200).JSON(utility.Message{Text: "Card successfully edited"})
}
