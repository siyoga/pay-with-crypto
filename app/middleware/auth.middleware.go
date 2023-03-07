package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func Auth(c *fiber.Ctx) error {
	accessKey := "secretAccessKey"
	// refreshKey := "secretRefreshKey"

	jwtware.New(jwtware.Config{
		SigningKey:   []byte(accessKey),
		ErrorHandler: jwtError,
	})
	c.Status(fiber.StatusOK)
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
