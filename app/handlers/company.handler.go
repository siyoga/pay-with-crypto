package handlers

import (
	"fmt"
	"path/filepath"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/datastore/s3"
	"pay-with-crypto/app/utility"

	"github.com/gofiber/fiber/v2"
)

func UpdateHandler(c *fiber.Ctx) error {
	var request utility.UpdateInfoRequest
	company := c.Locals("company").(db.Company)

	if err := c.BodyParser(&request); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request body"})
	}

	company.Image = request.Image
	company.LinkToCompany = request.Link
	company.Name = request.Name

	if !db.WholeOneUpdate(company) {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Failed to update entity on server"})
	}

	return c.Status(fiber.StatusOK).JSON(company)
}

func GetHandler(c *fiber.Ctx) error {
	var company db.Company
	var state bool

	companyId := c.Query("id")

	if companyId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request"})
	}

	company, state = db.GetOneByWithPreload[db.Company]("id", "Cards", companyId)

	if !state {
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Company not exist"})
	}

	return c.Status(fiber.StatusOK).JSON(company)
}

func CompanyLogoUploadHandler(c *fiber.Ctx) error {
	companyLogoRaw, err := c.FormFile("companyLogo")
	company := c.Locals("company").(db.Company)

	fmt.Println("requested")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request"})
	}

	companyLogo, err := companyLogoRaw.Open()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request, provide companyLogo correctly"})
	}

	defer companyLogo.Close()

	fileName := company.ID.String() + "_logo" + filepath.Ext(companyLogoRaw.Filename)
	s3.DeleteImage(fileName) // чтобы обезопаситься от дубликатов

	fileLink, isUploadOk := s3.UploadFile(companyLogo, fileName)

	if !isUploadOk {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	company, isUpdateOk := db.UpdateOneBy[db.Company]("id", company.ID, "image", fileLink)

	if !isUpdateOk {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something’s wrong with the server. Try it later."})
	}

	return c.Status(fiber.StatusOK).JSON(company)
}
