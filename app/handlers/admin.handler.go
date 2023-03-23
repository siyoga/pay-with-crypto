package handlers

import (
	"os"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/utility"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"

	"golang.org/x/crypto/bcrypt"
)

func AdminRegisterHandler(c *fiber.Ctx) error {
	var admin db.Admin

	if err := c.BodyParser(&admin); err != nil {
		return fiber.ErrBadRequest
	}

	if _, engaged := db.GetOneBy[db.Admin]("name", admin.Name); engaged {

		return fiber.ErrConflict
	}

	admin.ID = uuid.Must(uuid.NewV4())

	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 12)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	admin.Password = string(hash)

	if ok := db.Add(admin); !ok {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(admin)
}

func GetCardsForApprove(c *fiber.Ctx) error {
	cards, _ := db.GetManyBy[db.Card]("approved", "pending")

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

func CreateFirstAdmin() {
	var firstAdmin db.Admin

	firstAdmin.ID = uuid.Must(uuid.NewV4())
	firstAdmin.Name = os.Getenv("ADMIN_USERNAME")
	firstAdmin.FirstName = os.Getenv("ADMIN_FIRSTNAME")
	firstAdmin.LastName = os.Getenv("ADMIN_LASTNAME")
	firstAdmin.Password = os.Getenv("ADMIN_PASSWORD")
	hash, _ := bcrypt.GenerateFromPassword([]byte(firstAdmin.Password), 12)
	firstAdmin.Password = string(hash)

	if empty := db.AdminCheck(); empty {
		db.Add(firstAdmin)
	}
}

func AdminLoginHandler(c *fiber.Ctx) error {
	var requsetData db.Admin
	var refreshToken db.RefreshToken

	if err := c.BodyParser(&requsetData); err != nil {
		return fiber.ErrBadRequest
	}

	admin, state := db.Auth[db.Admin](requsetData.Name)
	if !state {
		return fiber.ErrBadRequest
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(requsetData.Password)); err != nil {
		return fiber.ErrBadRequest
	}

	response, errs := generatTokenResponse(admin.ID)
	if errs[0] != nil {
		return fiber.ErrInternalServerError
	}
	if errs[1] != nil {
		return fiber.ErrInternalServerError
	}

	refreshToken.Token = response.RefreshToken

	if ok := db.Add(refreshToken); !ok {
		return fiber.ErrBadRequest
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func ValidateCard(c *fiber.Ctx) error {
	var body utility.Status
	var response string

	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}

	if body.Status {
		db.UpdateOneBy[db.Card]("id", body.ID, "approved", "approved")
		response = "Card is approved"
	} else {
		db.UpdateOneBy[db.Card]("id", body.ID, "approved", "disapproved")
		response = "Card is disapproved"
	}

	return c.Status(200).JSON(response)
}

func SoftDeleteHandler(c *fiber.Ctx) error {
	var company db.Company
	var state bool

	if err := c.BodyParser(&company); err != nil {
		return fiber.ErrBadRequest
	}

	if company, state = db.GetOneBy[db.Company]("id", company.ID); !state {
		return fiber.ErrBadRequest
	}

	if state = db.DeleteBy[db.Company]("id", company.ID); !state {
		return fiber.ErrInternalServerError
	}

	return c.Status(200).JSON(fiber.Map{"message": "User deleted from scope"})
}
