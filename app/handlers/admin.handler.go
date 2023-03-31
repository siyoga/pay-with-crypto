package handlers

import (
	"fmt"
	"os"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/utility"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"

	"golang.org/x/crypto/bcrypt"
)

// @Description Create new admin accounts tokens.
// @Tags Auth
// @Accept json
// @Produce json
// @Param admin_data body object{username=string,first_name=string,last_name=string,password=string} true "Admin data"
// @Security ApiKeyAuth
// @Success 200 {object} datastore.Admin
// @Failure 409 {object} utility.Message "Admin already created"
// @Failure 400 {object} utility.Message "Invalid request body"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /auth/admin/register [post]
func AdminRegisterHandler(c *fiber.Ctx) error {
	var admin db.Admin

	if err := c.BodyParser(&admin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request body"})
	}

	if _, engaged := db.GetOneBy[db.Admin]("name", admin.Name); engaged {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Such an admin is already created"})
	}

	admin.ID = uuid.Must(uuid.NewV4())

	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}
	admin.Password = string(hash)

	if ok := db.Add(admin); !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.Status(fiber.StatusCreated).JSON(admin)
}

// @Description Get cards for validate
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []datastore.Card
// @Router /admin/getForApprove [get]
func GetCardsForApprove(c *fiber.Ctx) error {
	cards, _ := db.GetManyBy[db.Card]("approved", "pending")

	return c.Status(200).JSON(cards)
}

func CreateFirstAdmin() {
	var firstAdmin db.Admin

	firstAdmin.ID = uuid.Must(uuid.NewV4())
	firstAdmin.Name = os.Getenv("ADMIN_USERNAME")
	firstAdmin.FirstName = os.Getenv("ADMIN_FIRSTNAME")
	firstAdmin.LastName = os.Getenv("ADMIN_LASTNAME")
	firstAdmin.Password = os.Getenv("ADMIN_PASSWORD")
	hash, _ := bcrypt.GenerateFromPassword([]byte(firstAdmin.Password), 12)
	firstAdmin.Password = string(hash)

	if empty := db.AdminCheck(); empty {
		fmt.Println("Create default admin account...")
		db.Add(firstAdmin)
	}
}

// @Description Login to admin account
// @Tags Auth
// @Accept json
// @Produce json
// @Param admin_data body object{name=string,password=string} true "Admin data"
// @Success 200 {object} datastore.Admin
// @Failure 409 {object} utility.Message "Admin already created"
// @Failure 400 {object} utility.Message "Invalid request body"
// @Failure 400 {object} utility.Message "Invalid credentials"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /auth/admin/login [post]
func AdminLoginHandler(c *fiber.Ctx) error {
	var requestData db.Admin
	var refreshToken db.RefreshToken

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "Invalid request body"})
	}

	admin, state := db.Auth[db.Admin](requestData.Name)
	fmt.Println(admin)

	if !state {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(requestData.Password)); err != nil {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "Invalid credentials"})
	}

	response, errs := generateTokenResponse(admin.ID)
	if errs[0] != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}
	if errs[1] != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	refreshToken.Token = response.RefreshToken

	if ok := db.Add(refreshToken); !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// @Description Validate company card as admin
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param validate_data body object{id=string,status=bool} true "Validate data"
// @Success 200 {object} utility.Message
// @Failure 400 {object} utility.Message "Invalid request body"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /admin/validateCard [put]
func ValidateCard(c *fiber.Ctx) error {
	var body utility.Status
	var response string

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "Invalid request body"})
	}

	if body.Status {
		_, isOK := db.UpdateOneBy[db.Card]("id", body.ID, "approved", "approved")

		if !isOK {
			return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
		}

		response = "Card is approved"
	} else {
		_, isOK := db.UpdateOneBy[db.Card]("id", body.ID, "approved", "disapproved")

		if !isOK {
			return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
		}

		response = "Card is disapproved"
	}

	return c.Status(200).JSON(utility.Message{Text: response})
}

// @Description Ban company account
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param company_data body object{id=string} true "Company data"
// @Success 200 {object} utility.Message
// @Failure 400 {object} utility.Message "Invalid request body"
// @Failure 404 {object} utility.Message "Company not exist"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /admin/ban [delete]
func SoftDeleteHandler(c *fiber.Ctx) error {
	var company db.Company
	var state bool

	if err := c.BodyParser(&company); err != nil {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "Invalid request body"})
	}

	if company, state = db.GetOneBy[db.Company]("id", company.ID); !state {
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Such a company does not exist"})
	}

	if state = db.UnscopeCompanyByIdWithCards(company.ID); !state {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.Status(200).JSON(utility.Message{Text: "Company and card of this company deleted from scope."})
}

// @Description Unban company account
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param company_data body object{id=string} true "Company data"
// @Success 200 {object} utility.Message "Company added to server scope"
// @Failure 400 {object} utility.Message "Invalid request body"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /admin/unban [delete]
func UnbanCompanyHandler(c *fiber.Ctx) error {
	var company db.Company
	var state bool

	if err := c.BodyParser(&company); err != nil {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "Invalid request body"})
	}

	if state = db.ScopeCompanyByIdWithCards(company.ID); !state {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.Status(200).JSON(utility.Message{Text: "Company and card of this company added to scope."})

}

// @Description Validate tag as admin
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param validate_data body object{id=string,status=bool} true "Validate data"
// @Success 200 {object} utility.Message
// @Failure 400 {object} utility.Message "Invalid request body"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /admin/validateTag [put]
func ValidateTag(c *fiber.Ctx) error {
	var body utility.Status
	var response string

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusConflict).JSON(utility.Message{Text: "Invalid request body"})
	}

	if body.Status {
		_, isOK := db.UpdateOneBy[db.Tag]("id", body.ID, "approved", "approved")
		if !isOK {
			return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
		}

		response = "Tag is approved"
	} else {
		_, isOK := db.UpdateOneBy[db.Tag]("id", body.ID, "approved", "dispproved")

		if !isOK {
			return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
		}

		response = "Tag is disapproved"
	}

	return c.Status(200).JSON(utility.Message{Text: response})
}
