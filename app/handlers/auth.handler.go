package handlers

import (
	db "pay-with-crypto/app/datastore"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"

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

func LoginHandler(c *fiber.Ctx) error {
	var user db.User
	var loginResponse db.LoginResponse

	if err := c.BodyParser(&user); err != nil {
		return fiber.ErrBadRequest
	}

	user, state := db.UserAuth(user.Company_Name, user.Password)
	if !state {
		return fiber.ErrBadRequest
	}

	payload := jwt.MapClaims{
		"sub": user.ID,
	}

	token, err := db.GeneratToken("secretKey", payload)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	loginResponse.UserID = user.ID
	loginResponse.AccessToken = token

	return c.Status(fiber.StatusOK).JSON(loginResponse)
}
