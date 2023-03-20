package handlers

import (
	"fmt"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/utility"
	googleutil "pay-with-crypto/app/utility/google.util"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/grokify/go-pkce"
	"github.com/sethvargo/go-password/password"

	"net/mail"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

func RegisterHandler(c *fiber.Ctx) error {
	var user db.Company

	if err := c.BodyParser(&user); err != nil {
		return fiber.ErrBadRequest
	}

	if _, engaged := db.GetOneBy[db.Company]("name", user.Name); engaged {
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

func AuthGoogleGetApprove(c *fiber.Ctx) error {
	path := googleutil.ConfigGoogle()

	NewPKCE := *googleutil.CreatePKCE()

	url := path.AuthCodeURL("state",
		oauth2.SetAuthURLParam(pkce.ParamCodeChallenge, NewPKCE.CodeChallenge),
		oauth2.SetAuthURLParam(pkce.ParamCodeChallengeMethod, pkce.MethodS256)) //TODO!:CHANGE STATE TO SOMETHING MORE SECURITY STRONG. In theory it should be random each time, but not sure.
	return c.Redirect(url)
}

func Callback(c *fiber.Ctx) error {
	var refreshToken db.RefreshToken

	state := c.FormValue("state")
	if state != "state" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Missing or malformed URI"})
	}

	code := c.FormValue("code")
	PCKECode := googleutil.ThePKCE.CodeVerifier

	tokens, err := googleutil.GetTokens(code, PCKECode)
	if err != nil {
		fmt.Println(err.Error())
	}
	UserData, err := googleutil.GetUserData(tokens)
	if err != nil {
		return err
	}

	UserDataFromDb, check := db.GetOneBy[db.Company]("mail", UserData.Email)
	if check {
		response, err := AuthGoogleLoginUser(c, UserDataFromDb)
		if err != nil {
			return fiber.ErrBadRequest
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}

	_, err = mail.ParseAddress(UserData.Email)
	if err != nil {
		return fiber.ErrBadRequest
	}
	res, err := password.Generate(64, 10, 10, false, false)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	user := db.Company{
		ID:       uuid.Nil,
		Name:     UserData.Name,
		Image:    UserData.Picture,
		Password: res,
		Mail:     UserData.Email,
	}
	user.ID = uuid.Must(uuid.NewV4())

	hash, err := bcrypt.GenerateFromPassword([]byte(res), 12)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	user.Password = string(hash)

	if ok := db.Add(user); !ok {
		return fiber.ErrInternalServerError
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

func AuthGoogleLoginUser(c *fiber.Ctx, userdata db.Company) (utility.JWTTokenPair, error) {
	var refreshToken db.RefreshToken
	payload := jwt.MapClaims{
		"sub":       userdata.ID,
		"generated": time.Now().Add(15 * 24 * time.Hour),
	}

	response, errs := generatTokenResponse(payload)
	if errs[0] != nil {
		return utility.JWTTokenPair{}, errs[0]
	}
	if errs[1] != nil {
		return utility.JWTTokenPair{}, errs[1]
	}

	refreshToken.Token = response.RefreshToken

	if ok := db.Add(refreshToken); !ok {
		return utility.JWTTokenPair{}, fiber.ErrBadRequest
	}

	return response, nil
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
