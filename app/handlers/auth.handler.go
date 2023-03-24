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
	var requestData db.Company
	var refreshToken db.RefreshToken

	if err := c.BodyParser(&requestData); err != nil {
		return fiber.ErrBadRequest
	}

<<<<<<< HEAD
	user, state := db.Auth[db.Company](requestData.Name)
=======
	company, state := db.Auth[db.Company](requsetData.Name)
>>>>>>> 63fd6ba57a0daaf531e0867a58c619a75a7bc636
	if !state {
		return fiber.ErrBadRequest
	}

<<<<<<< HEAD
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestData.Password)); err != nil {
=======
	if err := bcrypt.CompareHashAndPassword([]byte(company.Password), []byte(requsetData.Password)); err != nil {
>>>>>>> 63fd6ba57a0daaf531e0867a58c619a75a7bc636
		return fiber.ErrBadRequest
	}

	response, errs := generatTokenResponse(company.ID)
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

func UpdateTokensHandler(c *fiber.Ctx) error {
	var refreshToken db.RefreshToken
	company := c.Locals("company").(db.Company)

	if err := c.BodyParser(&refreshToken); err != nil {
		return fiber.ErrBadRequest
	}

	if _, state := db.GetOneBy[db.RefreshToken]("token", refreshToken.Token); !state {
		return fiber.ErrBadRequest
	}

	response, errs := generatTokenResponse(company.ID)
	if errs[0] != nil {
		return fiber.ErrInternalServerError
	}
	if errs[1] != nil {
		return fiber.ErrInternalServerError
	}

	if _, ok := db.UpdateOneBy[db.RefreshToken]("token", string(refreshToken.Token), "token", string(response.RefreshToken)); !ok {
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

	response, errs := generatTokenResponse(user.ID)
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

	response, errs := generatTokenResponse(userdata.ID)
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

func generatTokenResponse(companyID uuid.UUID) (utility.JWTTokenPair, []error) {

	payload := jwt.MapClaims{
		"sub":       companyID,
		"generated": time.Now().Add(15 * 24 * time.Hour),
	}

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
