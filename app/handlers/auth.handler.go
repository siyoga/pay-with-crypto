package handlers

import (
	db "pay-with-crypto/app/datastore"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *fiber.Ctx) error {
	var user db.User

	if err := c.BodyParser(&user); err != nil {
		return fiber.ErrBadRequest
	}

	user.ID = uuid.Must(uuid.NewV4())

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	user.Password = string(hash)

	if ok := db.Add(user); !ok {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
