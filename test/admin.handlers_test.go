package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"net/http/httptest"

	"os"
	"testing"

	"pay-with-crypto/app/controllers"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/utility"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// The tests require existence of the first admin

func TestGetCardsForApprove(t *testing.T) {
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
	controllers.AdminController(testServer)

	// connect db
	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	firstAdmin := db.Admin{
		Name:     "admin",
		Password: "1234567890"}

	tokens := adminLogin(firstAdmin, testServer)

	tests := []struct {
		name       string
		method     string
		path       string
		statusCode int
	}{
		{
			name:       "Successful pipeline",
			method:     http.MethodGet,
			path:       "http://127.0.0.1:8081/admin/getForApprove",
			statusCode: 200,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testRequest := httptest.NewRequest(tc.method, tc.path, nil)
			testRequest.Header.Add("accessToken", tokens.AccessToken)

			response, _ := testServer.Test(testRequest)

			receivedStatusCode := response.StatusCode

			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}
		})
	}
}

func TestCreateTag(t *testing.T) {
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
	controllers.AdminController(testServer)

	// connect db
	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	firstAdmin := db.Admin{
		Name:     "admin",
		Password: "1234567890"}

	tokens := adminLogin(firstAdmin, testServer)

	tests := []struct {
		name        string
		method      string
		path        string
		requestBody db.Tag
		statusCode  int
	}{
		{
			name:        "Successful pipeline",
			method:      http.MethodPost,
			path:        "http://127.0.0.1:8081/admin/createTag",
			requestBody: db.Tag{Name: "TEST_TOKEN"},
			statusCode:  201,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsedRequestBody, _ := json.Marshal(tc.requestBody)

			testRequest := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(parsedRequestBody))
			testRequest.Header.Add("Content-Type", "application/json")
			testRequest.Header.Add("accessToken", tokens.AccessToken)

			response, _ := testServer.Test(testRequest)

			receivedStatusCode := response.StatusCode

			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}
		})
	}

	db.DeleteBy[db.Tag]("name", "TEST_TOKEN")
}

func TestValidateCard(t *testing.T) {
	// Works even if the card does not exist
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
	controllers.AdminController(testServer)

	// connect db
	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	firstAdmin := db.Admin{
		Name:     "admin",
		Password: "1234567890"}

	tokens := adminLogin(firstAdmin, testServer)
	var id uuid.UUID

	tests := []struct {
		name        string
		method      string
		path        string
		requestBody utility.Status
		statusCode  int
	}{
		{
			name:        "Approve",
			method:      http.MethodPut,
			path:        "http://127.0.0.1:8081/admin/validateCard",
			requestBody: utility.Status{ID: id, Status: true},
			statusCode:  200,
		},
		{
			name:        "Disapprove",
			method:      http.MethodPut,
			path:        "http://127.0.0.1:8081/admin/validateCard",
			requestBody: utility.Status{ID: id, Status: false},
			statusCode:  200,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsedRequestBody, _ := json.Marshal(tc.requestBody)

			testRequest := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(parsedRequestBody))
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

func TestSoftDelete(t *testing.T) {
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
	controllers.AdminController(testServer)

	// connect db
	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	firstAdmin := db.Admin{
		Name:     "admin",
		Password: "1234567890"}

	tokens := adminLogin(firstAdmin, testServer)

	var company db.Company
	company.ID = uuid.Must(uuid.NewV4())
	company.Name = "SOFT_DELETED_COMPANY"
	company.Password = "password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(company.Password), 12)
	company.Password = string(hash)
	db.Add(company)

	tests := []struct {
		name        string
		method      string
		path        string
		requestBody db.Company
		statusCode  int
	}{
		{
			name:        "Successful pipeline",
			method:      http.MethodDelete,
			path:        "http://127.0.0.1:8081/admin/softDelete",
			requestBody: db.Company{ID: company.ID},
			statusCode:  200,
		},
		{
			name:        "The company is already soft deleted",
			method:      http.MethodDelete,
			path:        "http://127.0.0.1:8081/admin/softDelete",
			requestBody: db.Company{ID: company.ID},
			statusCode:  400,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsedRequestBody, _ := json.Marshal(tc.requestBody)

			testRequest := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(parsedRequestBody))
			testRequest.Header.Add("Content-Type", "application/json")
			testRequest.Header.Add("accessToken", tokens.AccessToken)

			response, _ := testServer.Test(testRequest)

			receivedStatusCode := response.StatusCode

			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}
		})
	}
	// The next line does not work, although it should
	db.DeleteBy[db.Company]("Name", "SOFT_DELETED_COMPANY")
}
