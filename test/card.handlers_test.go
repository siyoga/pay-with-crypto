package test

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
	"github.com/lib/pq"
)

var testCardResponse db.Card

var testCard = db.Card{
	Name:        "TestName",
	LinkToProd:  "google.com",
	Description: "TestDescription",
	Tags:        pq.StringArray([]string{"testtag1", "testtag2"}),
	Price:       "testPrice",
}

var testCompany = db.Company{
	Name:          "TestCreateCard",
	Password:      "TestCreateCardById_PASSWORD",
	Mail:          "testtt@kekw.com",
	LinkToCompany: "http://linktt.com",
}

func TestCreateCard(t *testing.T) {
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
	controllers.CardController(testServer)

	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	tests := []struct {
		name       string
		method     string
		path       string
		card       db.Card
		statusCode int
	}{
		{
			name:       "Successfull pipeline",
			method:     http.MethodPost,
			path:       "/card/create",
			card:       testCard,
			statusCode: 201,
		},
		{
			name:       "Exist card. Conflict",
			method:     http.MethodPost,
			path:       "/card/create",
			card:       testCard,
			statusCode: 409,
		},
	}

	auth(testCompany, testServer)
	tokens := login(db.Company{Name: testCompany.Name, Password: testCompany.Password}, testServer)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encodedJson, _ := json.Marshal(tc.card)
			testRequest := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(encodedJson))
			testRequest.Header.Add("Content-Type", "application/json")
			testRequest.Header.Add("accessToken", tokens.AccessToken)

			response, _ := testServer.Test(testRequest)
			body, _ := io.ReadAll(response.Body)

			var card db.Card
			_ = json.Unmarshal(body, &tokens)
			testCardResponse = card
			receivedStatusCode := response.StatusCode
			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}
		})
	}
}

func TestEditCard(t *testing.T) {
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
	controllers.CardController(testServer)

	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	tests := []struct {
		name       string
		method     string
		path       string
		card       db.Card
		statusCode int
	}{
		{
			name:       "Successfull pipeline",
			method:     http.MethodPost,
			path:       "/card/create",
			card:       db.Card{ID: testCardResponse.ID, Name: "newtestname"},
			statusCode: 200,
		},
		{
			name:       "Without id",
			method:     http.MethodPost,
			path:       "/card/create",
			card:       db.Card{Name: "newtestname"},
			statusCode: 400,
		},
	}

	tokens := login(db.Company{Name: testCompany.Name, Password: testCompany.Password}, testServer)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encodedJson, _ := json.Marshal(tc.card)
			testRequest := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(encodedJson))
			testRequest.Header.Add("Content-Type", "application/json")
			testRequest.Header.Add("accessToken", tokens.AccessToken)

			response, _ := testServer.Test(testRequest)
			receivedStatusCode := response.StatusCode
			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}
		})
	}
}

func TestSearchIdCard(t *testing.T) {
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
	controllers.CardController(testServer)

	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	tests := []struct {
		name       string
		method     string
		path       string
		card       db.Card
		statusCode int
	}{
		{
			name:       "Successfull pipeline",
			method:     http.MethodPost,
			path:       "/card/create?id=" + testCardResponse.ID.String(),
			statusCode: 200,
		},
		{
			name:       "Without id",
			method:     http.MethodPost,
			path:       "/card/create?id=",
			statusCode: 400,
		},
	}

	tokens := login(db.Company{Name: testCompany.Name, Password: testCompany.Password}, testServer)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encodedJson, _ := json.Marshal(tc.card)
			testRequest := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(encodedJson))
			testRequest.Header.Add("Content-Type", "application/json")
			testRequest.Header.Add("accessToken", tokens.AccessToken)

			response, _ := testServer.Test(testRequest)
			receivedStatusCode := response.StatusCode
			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}
		})
	}
}
