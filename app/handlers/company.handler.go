package handlers

import (
	"os"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/datastore/s3"
	"pay-with-crypto/app/utility"

	"github.com/gofiber/fiber/v2"
)

// @Description Search company by id
// @Tags Company
// @Accept json
// @Produce json
// @Param companyId query string true "Company id"
// @Success 200 {object} datastore.Company
// @Failure 400 {object} utility.Message "Invalid request"
// @Failure 404 {object} utility.Message "No card"
// @Failure 403 {object} utility.Message "Card owner was banned"
// @Router /company/search/id [get]
func CompanyGetByIdHandler(c *fiber.Ctx) error {
	var company db.Company
	var state bool

	companyId := c.Query("id")

	if companyId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request"})
	}

	company, state = db.GetUserById(companyId)

	if !state {
		if _, state = db.GetOneUnscopedBy[db.Card]("id", companyId); state {
			return c.Status(fiber.StatusForbidden).JSON(utility.Message{Text: "Owner of card was banned"})
		}
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Company not exist"})
	}

	return c.Status(fiber.StatusOK).JSON(company)
}

// @Description Company logo uploader
// @Tags Company
// @Accept json
// @Produce json
// @Security accessToken
// @Param companyLogo formData file true "Logo image"
// @Success 204 "Company logo successful uploaded"
// @Failure 400 {object} utility.Message "Invalid request, log in"
// @Failure 400 {object} utility.Message "Invalid request, provide companyLogo"
// @Failure 500 {object} utility.Message "Internal server error"
// @Router /company/uploadLogo [post]
func CompanyLogoUploaderHandler(c *fiber.Ctx) error {
	logoBucket := os.Getenv("S3_BUCKET_COMPANY_LOGOS")
	companyLogo, err := c.FormFile("companyLogo")
	company := c.Locals("company").(db.Company)

	companyId := company.ID.String()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request, log in to account"})
	}

	companyLogoBuffer, err := companyLogo.Open()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request, provide companyLogo"})
	}

	defer companyLogoBuffer.Close()

	fileName := companyId + "_logo"
	fileNameInS3, isUploadOk := s3.UploadFile(companyLogo, companyLogoBuffer, logoBucket, fileName)

	if !isUploadOk {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	_, isUpdateOk := db.UpdateOneBy[db.Company]("id", companyId, "image", *fileNameInS3)

	if !isUpdateOk {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
