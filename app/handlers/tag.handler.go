package handlers

import (
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/utility"

	"github.com/gofrs/uuid"

	"github.com/gofiber/fiber/v2"
)

// @Description Create tag as company
// @Tags Tag
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param tag_data body object{name=string} true "Tag data"
// @Success 201 "Tag created"
// @Failure 400 {object} utility.Message "Invalid request"
// @Failure 409 {object} utility.Message "Tag already exist"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /tag/create [post]
func TagCreateHandler(c *fiber.Ctx) error {
	var newTag db.Tag
	company := c.Locals("company").(db.Company)

	if err := c.BodyParser(&newTag); err != nil {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "Invalid request body"})
	}

	if _, isExist := db.GetOneBy[db.Tag]("name", newTag.Name); isExist {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "This tag is already exist"})
	}

	newTag.ID = uuid.Must(uuid.NewV4())
	newTag.CreatorID = company.ID
	newTag.Approved = "pending"

	if ok := db.Add(newTag); !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.Status(fiber.StatusCreated).JSON(newTag)
}

// @Description Return all approved tags
// @Tags Tag
// @Accept json
// @Produce json
// @Success 200 "Return tags"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /tag/get/all [get]
func TagGetterHandler(c *fiber.Ctx) error {
	tags, isOk := db.GetManyBy[db.Tag]("approved", "approved")

	if !isOk {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.Status(fiber.StatusOK).JSON(tags)
}
