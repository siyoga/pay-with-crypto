package handlers

import (
	db "pay-with-crypto/app/datastore"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"

	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *fiber.Ctx) error {
	var user db.User

	if err := c.BodyParser(&user); err != nil {
		return fiber.ErrBadRequest
	}

	_, err := mail.ParseAddress(user.Mail)
	if err != nil {
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
	var requsetData db.User
	var refreshToken db.RefreshToken

	if err := c.BodyParser(&requsetData); err != nil {
		return fiber.ErrBadRequest
	}

	user, state := db.UserAuth(requsetData.Company_Name, requsetData.Password)
	if !state {
		return fiber.ErrBadRequest
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requsetData.Password)); err != nil {
		return fiber.ErrBadRequest
	}

	payload := jwt.MapClaims{
		"sub":       user.ID,
		"generated": time.Now().Add(15 * 24 * time.Hour),
	}

	response, errs := generatTokenResponse(payload)
	if errs != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	refreshToken.Token = response[1]

	return c.Status(fiber.StatusOK).SendString(response[0])
}

func generatTokenResponse(payload jwt.MapClaims) ([]string, []error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	accessToken, err_access := token.SignedString("secretAccessKey")
	refreshToken, err_refresh := token.SignedString("secretRefreshKey")

	response := make([]string, 2)
	response[0] = accessToken
	response[1] = refreshToken

	errors := make([]error, 2)
	errors[0] = err_access
	errors[1] = err_refresh

	return response, errors

}
