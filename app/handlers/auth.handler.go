package handlers

import (
	"fmt"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/utility"
	google "pay-with-crypto/app/utility/google"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *fiber.Ctx) error {
	var request db.Company
	var refreshToken db.RefreshToken

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid request body"})
	}

	company, state := db.Auth[db.Company](request.Email)

	if !state {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid credentials"})
	}

	if !company.ViaGoogle {
		if err := bcrypt.CompareHashAndPassword([]byte(company.Password), []byte(request.Password)); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid credentials"})
		}
	}

	tokens := generateTokens(company.ID)

	refreshToken.Token = tokens.RefreshToken
	refreshToken.CompanyID = company.ID

	if ok := db.Add(refreshToken); !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something went wrong on our side"})
	}

	return c.Status(fiber.StatusOK).JSON(tokens)
}

func RefreshHandler(c *fiber.Ctx) error {
	var refreshToken db.RefreshToken

	if err := c.BodyParser(&refreshToken); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Refresh token was not provided"})
	}

	if _, state := db.GetOneBy[db.RefreshToken]("token", refreshToken.Token); !state {
		return c.Status(fiber.StatusNotFound).JSON(utility.Message{Text: "Such refresh token is not exist"})
	}

	userId, err := unwrapRefreshJWT(refreshToken.Token)

	if !err {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Invalid refresh token, try log in"})
	}

	response := generateTokens(uuid.FromStringOrNil(userId))

	if _, ok := db.UpdateOneBy[db.RefreshToken]("token", string(refreshToken.Token), "token", string(response.RefreshToken)); !ok {
		return c.Status(fiber.StatusBadRequest).JSON(utility.Message{Text: "Can't update refresh token"})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func WhoAmIHandler(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	atSecretKey := utility.GetEnv("ACCESS_SECRET", "secretAccessKey")
	accessToken := strings.Split(authHeader, "Bearer ")[1]
	claims := &utility.Claims{}

	accessTokenPayload, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(atSecretKey), nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something went wrong on our side"})
	}

	if !accessTokenPayload.Valid {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Invalid token"})
	}

	company, ok := db.GetOneBy[db.Company]("id", claims.Sub)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something went wrong on our side"})
	}

	return c.JSON(company)
}

func RegisterHandler(c *fiber.Ctx) error {
	var registerInfoRequest utility.RegisterInfoRequest
	var company db.Company
	var gToken string

	if err := c.BodyParser(&registerInfoRequest); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utility.Message{Text: "Invalid register payload"})
	}

	company.Email = registerInfoRequest.Email

	existUser, isOk := db.GetOneBy[db.Company]("email", company.Email)

	if isOk {
		return c.Status(fiber.StatusOK).JSON(existUser)
	}

	company.ID = uuid.Must(uuid.NewV4())

	company.CreatedAt = time.Now()
	company.UpdateAt = time.Now()

	// Здесь добавляем другие провайдеры, в том числе и простой вход через почту.
	if registerInfoRequest.ViaGoogle {
		company.ViaGoogle = true

		gToken = c.Get("Authorization", "")
		if gToken == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utility.Message{Text: "Invalid id_token"})
		}

		_, gErr, err := google.GetInfoByIdToken(gToken)
		if gErr != (google.GoogleErrorResponse{}) {
			return c.Status(fiber.StatusUnauthorized).JSON(utility.Message{Text: "Invalid id_token"})
		}

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something went wrong on our side"})
		}

		password := uuid.Must(uuid.NewV4()).String()
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something went wrong on our side"})
		}

		company.Password = string(hash)
	}

	if ok := db.Add(company); !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(utility.Message{Text: "Something went wrong on our side"})
	}

	return c.Status(fiber.StatusCreated).JSON(company)
}

// ---------------------PRIVATE------------------------

func generateTokens(ID uuid.UUID) utility.JWTTokenPair {
	var tokens utility.JWTTokenPair

	atPayload := jwt.MapClaims{
		"sub":       ID,
		"generated": time.Now(),
		"dead":      time.Now().Add(24 * time.Hour),
	}

	rtPayload := jwt.MapClaims{
		"sub":       ID,
		"generated": time.Now(),
		"dead":      time.Now().Add(7 * 24 * time.Hour),
	}

	atSecretKey := utility.GetEnv("ACCESS_SECRET", "secretAccessKey")
	rtSecretKey := utility.GetEnv("REFRESH_SECRET", "secretRefreshKey")

	fmt.Print(atSecretKey, rtSecretKey)

	accessToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, atPayload).SignedString([]byte(atSecretKey))
	refreshToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, rtPayload).SignedString([]byte(rtSecretKey))

	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return tokens
}

func unwrapRefreshJWT(tokenString string) (string, bool) {
	var userId string
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secretRefreshKey"), nil
	})

	if err != nil {
		return "", false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId = claims["sub"].(string)
	} else {
		return "", false
	}

	return userId, true
}
