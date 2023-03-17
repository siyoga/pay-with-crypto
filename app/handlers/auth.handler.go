package handlers

import (
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/utility"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"

	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c *fiber.Ctx) error {
	var user db.Company

	if err := c.BodyParser(&user); err != nil {
		return fiber.ErrBadRequest
	}

	if empty := db.ExportedIsUniqueCompany(user.Name); !empty {

		return fiber.ErrConflict
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
	var requsetData db.Company
	var refreshToken db.RefreshToken

	if err := c.BodyParser(&requsetData); err != nil {
		return fiber.ErrBadRequest
	}

	user, state := db.Auth[db.Company](requsetData.Name)
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

func generatTokenResponse(payload jwt.MapClaims) (utility.JWTTokenPair, []error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	var response utility.JWTTokenPair
	errors := make([]error, 2)

	accessToken, err_access := token.SignedString([]byte("secretAccessKey"))
	refreshToken, err_refresh := token.SignedString([]byte("secretRefreshKey"))

	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	errors[0] = err_access
	errors[1] = err_refresh

	return response, errors

}
