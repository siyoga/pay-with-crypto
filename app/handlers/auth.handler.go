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

var internalServerError = utility.ErrorResponse{Message: "Such a company is already registered."}

// @Description Create company.
// @Tags Auth
// @Accept json
// @Produce json
// @Param name body string true "Company name"
// @Param password body string true "Company password"
// @Param mail body string true "Company mail"
// @Param linkToCompany body string true "Link to company website"
// @Success 201 {object} datastore.Company
// @Failure 409 {object} utility.ErrorMessage "Company already created"
// @Failure 400 {object} utility.ErrorMessage "Invalid company email"
// @Failure 500 {object} utility.ErrorMessage "Internal server error"
// @Router /auth/register [post]
func RegisterHandler(c *fiber.Ctx) error {
	var company db.Company

	if err := c.BodyParser(&company); err != nil {
		return fiber.ErrBadRequest
	}

	if _, exist := db.GetOneBy[db.Company]("name", company.Name); exist {
		return c.Status(fiber.StatusConflict).JSON(internalServerError)
	}

	_, err := mail.ParseAddress(company.Mail)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.ErrorResponse{Message: "Enter a valid email address"})
	}

	company.ID = uuid.Must(uuid.NewV4())

	hash, err := bcrypt.GenerateFromPassword([]byte(company.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(internalServerError)
	}
	company.Password = string(hash)

	if ok := db.Add(company); !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(internalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(company)
}

// @Description Login to company account.
// @Tags Auth
// @Accept json
// @Produce json
// @Param name body string true "Company name"
// @Param password body string true "Company password"
// @Success 200 {object} utility.JWTTokenPair
// @Failure 409 {object} utility.ErrorMessage "Company already created"
// @Failure 400 {object} utility.ErrorMessage "Invalid company email"
// @Failure 500 {object} utility.ErrorMessage "Internal server error"
// @Router /auth/login [post]
func LoginHandler(c *fiber.Ctx) error {
	var requestData db.Company
	var refreshToken db.RefreshToken

	if err := c.BodyParser(&requestData); err != nil {
		return fiber.ErrBadRequest
	}

	company, state := db.Auth[db.Company](requestData.Name)

	if !state {
		return c.Status(fiber.StatusBadRequest).JSON(utility.ErrorResponse{Message: "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(company.Password), []byte(requestData.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.ErrorResponse{Message: "Invalid credentials"})
	}

	response, errs := generateTokenResponse(company.ID)
	if errs[0] != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(internalServerError)
	}
	if errs[1] != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(internalServerError)
	}

	refreshToken.Token = response.RefreshToken

	if ok := db.Add(refreshToken); !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(internalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// @Description Update tokens.
// @Tags Auth
// @Accept json
// @Produce json
// @Param token body string true "Refresh token"
// @Success 200 {object} utility.JWTTokenPair
// @Failure 409 {object} utility.ErrorMessage "Token already created"
// @Failure 400 {object} utility.ErrorMessage "Refresh token was not provided"
// @Failure 400 {object} utility.ErrorMessage "Can't update refresh token"
// @Failure 500 {object} utility.ErrorMessage "Internal server error"
// @Router /auth/token_update [post]
func UpdateTokensHandler(c *fiber.Ctx) error {
	var refreshToken db.RefreshToken
	company := c.Locals("company").(db.Company)

	if err := c.BodyParser(&refreshToken); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Refresh token was not provided")
	}

	if _, state := db.GetOneBy[db.RefreshToken]("token", refreshToken.Token); state {
		return c.Status(fiber.StatusConflict).JSON("Such refresh token already exist")
	}

	response, errs := generateTokenResponse(company.ID)
	if errs[0] != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(internalServerError)
	}
	if errs[1] != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(internalServerError)
	}

	if _, ok := db.UpdateOneBy[db.RefreshToken]("token", string(refreshToken.Token), "token", string(response.RefreshToken)); !ok {
		return c.Status(fiber.StatusBadRequest).JSON("Can't update refresh token")
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

	response, errs := generateTokenResponse(user.ID)
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

	response, errs := generateTokenResponse(userdata.ID)
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

func generateTokenResponse(companyID uuid.UUID) (utility.JWTTokenPair, []error) {

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
