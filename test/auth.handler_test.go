package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"pay-with-crypto/app/controllers"
	db "pay-with-crypto/app/datastore"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		requestBody db.User
		statusCode  int
	}{
		{
			name:        "Successful pipeline",
			method:      http.MethodPost,
			path:        "/auth/register",
			statusCode:  201,
			requestBody: db.User{Company_Name: "TEST_COMPANY_NAME", Password: "TEST_PASSWORD"},
		},
	}

	if err := godotenv.Load("../dev.env"); err != nil {
		log.Fatalf("Error loading dev.env file")
	}

	config := db.DatabaseConfig{
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
		Host:     os.Getenv("DATABASE_HOST"),
	}

	testServer := fiber.New()
	controllers.AuthController(testServer)

	// connect db
	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsedRequestBody, _ := json.Marshal(tc.requestBody)
			testRequest := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(parsedRequestBody))
			testRequest.Header.Add("Content-Type", "application/json")

			response, _ := testServer.Test(testRequest)

			receivedStatusCode := response.StatusCode
			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}

			body, err := io.ReadAll(response.Body)

			if err != nil {
				t.Error("Incorrect response body, check the sent json")
			}

			var createdTestUser db.User

			if err := json.Unmarshal(body, &createdTestUser); err != nil {
				t.Errorf("Can't convert to User struct. Error description: %s", err.Error())
			}

			// TODO: Add crypt check
			// if createdTestUser.Password == tc.requestBody.Password {
			// 	t.Errorf("The user's password needs to be encrypted. Received: %s", createdTestUser.Password)
			// }
		})
	}
}
