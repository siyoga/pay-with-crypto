package handlers

import (
	"os"
	db "pay-with-crypto/app/datastore"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"

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
