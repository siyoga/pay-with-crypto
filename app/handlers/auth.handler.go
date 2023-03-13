package handlers

import (
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/utility"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"honnef.co/go/tools/config"

	"net/mail"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

func AuthGoogle(c *fiber.Ctx) error {
	path := &oauth2.Config{
		ClientID:     config.Config("GOOGLE_CLIENT"),
		ClientSecret: config.Config("GOOGLE_SECRET"),
		RedirectURL:  config.Config("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email"}, // you can use other scopes to get more data
		Endpoint: google.Endpoint,
	}
	url := path.AuthCodeURL("state")
	return c.Redirect(url)
}
