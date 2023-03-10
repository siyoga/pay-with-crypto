package handlers

import (
	"os"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/datastore/s3"

	"github.com/gofiber/fiber/v2"
)

func CompanyGetByIdHandler(c *fiber.Ctx) error {
	var company db.User
	var state bool

	companyId := c.Query("id")

	if companyId == "" {
		return fiber.ErrBadRequest
	}

	company, state = db.GetUserById(companyId)

	if !state {
		return fiber.ErrNotFound
	}

	return c.Status(fiber.StatusOK).JSON(company)
}

func CompanyLogoUploaderHandler(c *fiber.Ctx) error {
	logoBucket := os.Getenv("S3_BUCKET_COMPANY_LOGOS")
	companyLogo, err := c.FormFile("companyLogo")
	companyId := c.Query("companyId")

	if companyId == "" {
		return fiber.ErrBadRequest
	}

	if err != nil {
		return fiber.ErrBadRequest
	}

	cardLogoBuffer, err := companyLogo.Open()

	if err != nil {
		return fiber.ErrBadRequest
	}

	defer cardLogoBuffer.Close()

	fileName := companyId + "_logo"
	fileNameInS3, isUploadOk := s3.UploadFile(companyLogo, cardLogoBuffer, logoBucket, fileName)

	if !isUploadOk {
		return fiber.ErrInternalServerError
	}

	_, isUpdateOk := db.UpdateOneBy[db.User]("id", companyId, "image", *fileNameInS3)

	if !isUpdateOk {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusNoContent)
}
