package handlers

import (
	"os"
	db "pay-with-crypto/app/datastore"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"

	"golang.org/x/crypto/bcrypt"
)

func AdminRegisterHandler(c *fiber.Ctx) error {
	var admin db.Admin

	if err := c.BodyParser(&admin); err != nil {
		return fiber.ErrBadRequest
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

func CreateFirstAdmin() {
	var firstAdmin db.Admin

	firstAdmin.ID = uuid.Must(uuid.NewV4())
	firstAdmin.UserName = os.Getenv("ADMIN_USERNAME")
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

	admin, state := db.AdminAuth(requsetData.UserName, requsetData.Password)
	if !state {
		return fiber.ErrBadRequest
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(requsetData.Password)); err != nil {
		return fiber.ErrBadRequest
	}

	payload := jwt.MapClaims{
		"sub":       admin.ID,
		"generated": time.Now().Add(15 * 24 * time.Hour),
	}

	response, errs := generatTokenResponse(payload)
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
