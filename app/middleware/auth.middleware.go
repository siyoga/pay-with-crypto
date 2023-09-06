package middleware

import (
	"fmt"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/utility"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Auth(c *fiber.Ctx) error {
	var userid string

	accessKey := utility.GetEnv("ACCESS_SECRET", "secretAccessKey")

	authHeader := strings.Split(c.Get("Authorization"), " ")
	if len(authHeader) <= 1 {
		return c.Status(fiber.StatusUnauthorized).JSON(utility.Message{Text: "Invalid token format"})
	}

	accessToken := authHeader[1]

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(accessKey), nil
	})

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: err.Error()})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userid = claims["sub"].(string)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Missing or malformed JWT"})
	}

	result, ok := db.GetOneBy[db.Company]("id", userid)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Not found user"})
	}

	c.Locals("company", result)
	return c.Next()
}

func AuthAdmin(c *fiber.Ctx) error {
	var userid string

	hmacSampleSecret := "secretAccessKey"

	tokenString := c.Get("accessToken")
	if tokenString == "" {
		return fiber.ErrUnauthorized
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(hmacSampleSecret), nil
	})

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userid = claims["sub"].(string)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Missing or malformed JWT"})
	}

	result, ok := db.GetOneBy[db.Admin]("id", userid)
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Not found user"})
	}

	c.Locals("admin", result)
	return c.Next()
}
