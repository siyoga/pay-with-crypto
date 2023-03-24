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
	"pay-with-crypto/app/handlers"
	"pay-with-crypto/app/utility"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// Creating a global value to use tokens of this test session in other tests=

func TestRegisterHandler(t *testing.T) {
	if err := godotenv.Load("../dev.env"); err != nil {
		log.Fatalf("Error loading dev.env file")
	}

	testServer := fiber.New()
	controllers.AuthController(testServer)

	config := db.DatabaseConfig{
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
		Host:     os.Getenv("DATABASE_HOST"),
	}

	testCompany := db.Company{
		Name:          "TEST_COMPANY_NAME",
		Password:      "TEST_COMPANY_NAME_PASSWORD",
		Mail:          "testtt@kekw.com",
		LinkToCompany: "http://linktt.com"}

	tests := []struct {
		name        string
		method      string
		path        string
		requestBody db.Company
		statusCode  int
	}{
		{
			name:       "Successful company registeration",
			method:     http.MethodPost,
			path:       "/auth/register",
			statusCode: 201,
			requestBody: db.Company{Name: testCompany.Name,
				Password:      testCompany.Password,
				Mail:          testCompany.Mail,
				LinkToCompany: testCompany.LinkToCompany},
		},
		{
			name:        "Not enough data",
			method:      http.MethodPost,
			path:        "/auth/register",
			statusCode:  400,
			requestBody: db.Company{},
		},

		{
			name:       "Company with this name already exist",
			method:     http.MethodPost,
			path:       "/auth/register",
			statusCode: 409,
			requestBody: db.Company{Name: testCompany.Name,
				Password:      testCompany.Password,
				Mail:          testCompany.Mail,
				LinkToCompany: testCompany.LinkToCompany},
		},
		{
			name:       "Bad email",
			method:     http.MethodPost,
			path:       "/auth/register",
			statusCode: 400,
			requestBody: db.Company{Name: "test_name",
				Password:      "test_password",
				Mail:          "tsrrjfhmgduykuy",
				LinkToCompany: "test_link"},
		},
	}

	// connect db
	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	//handlers.CreateFirstAdmin()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsedRequestBody, _ := json.Marshal(tc.requestBody)
			testRequest := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(parsedRequestBody))
			testRequest.Header.Add("Content-Type", "application/json")

			fmt.Println(tc.requestBody)

			response, _ := testServer.Test(testRequest)

			receivedStatusCode := response.StatusCode
			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}
		})
	}

	//deleting data from test database
	db.DeleteBy[db.Company]("Name", testCompany.Name)

}

func TestLoginHandler(t *testing.T) {
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

	testCompany := db.Company{
		Name:          "TEST_COMPANY_NAME",
		Password:      "TEST_COMPANY_NAME_PASSWORD",
		Mail:          "testtt@kekw.com",
		LinkToCompany: "http://linktt.com"}

	company := auth(testCompany, testServer)

	tests := []struct {
		name        string
		method      string
		path        string
		requestBody db.Company
		statusCode  int
	}{
		{
			name:        "Successful company login",
			method:      http.MethodPost,
			path:        "/auth/login",
			statusCode:  200,
			requestBody: db.Company{Name: testCompany.Name, Password: testCompany.Password},
		},
		{
			name:        "Not enought data for login",
			method:      http.MethodPost,
			path:        "/auth/login",
			statusCode:  400,
			requestBody: db.Company{},
		},
		{
			name:        "Bad login",
			method:      http.MethodPost,
			path:        "/auth/login",
			statusCode:  400,
			requestBody: db.Company{Name: "Bad Name", Password: testCompany.Password},
		},
		{
			name:        "Bad password",
			method:      http.MethodPost,
			path:        "/auth/login",
			statusCode:  400,
			requestBody: db.Company{Name: testCompany.Name, Password: "Bad Password"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			parsedRequestBody, _ := json.Marshal(tc.requestBody)
			testRequest := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(parsedRequestBody))
			testRequest.Header.Add("Content-Type", "application/json")

			fmt.Println(tc.requestBody)

			response, _ := testServer.Test(testRequest)

			receivedStatusCode := response.StatusCode
			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}

		})
	}
	//deleting data from test database
	db.DeleteBy[db.Company]("Name", company.Name)

}

func TestUpdateTokensHandler(t *testing.T) {
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

	testCompany := db.Company{
		Name:          "TEST_COMPANY_NAME",
		Password:      "TEST_COMPANY_NAME_PASSWORD",
		Mail:          "testtt@kekw.com",
		LinkToCompany: "http://linktt.com"}

	company := auth(testCompany, testServer)

	tokens := login(db.Company{Name: testCompany.Name, Password: testCompany.Password}, testServer)

	tests := []struct {
		name        string
		method      string
		path        string
		requestBody db.RefreshToken
		statusCode  int
	}{
		{
			name:        "Successful tokens update",
			method:      http.MethodPost,
			path:        "/auth/token_update",
			statusCode:  200,
			requestBody: db.RefreshToken{Token: tokens.RefreshToken},
		},
		{
			name:        "Unexpected refresh token",
			method:      http.MethodPost,
			path:        "/auth/token_update",
			statusCode:  400,
			requestBody: db.RefreshToken{Token: "hjkluESDDLKJHZDFGB;'OJADFBLJIQEWFp'jio"},
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

	//deleting data from test database
	db.DeleteBy[db.Company]("Name", company.Name)

}

func TestAdminRegisterHandler(t *testing.T) {
	if err := godotenv.Load("../dev.env"); err != nil {
		log.Fatalf("Error loading dev.env file")
	}

	testServer := fiber.New()
	controllers.AuthController(testServer)

	config := db.DatabaseConfig{
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
		Host:     os.Getenv("DATABASE_HOST"),
	}

	testAdmin := db.Admin{Name: "greatestAdmin13378",
		FirstName: "Alexander",
		LastName:  "Nevsky",
		Password:  "superpassword"}

	tests := []struct {
		name        string
		method      string
		path        string
		requestBody interface{}
		statusCode  int
	}{
		{
			name:       "Successful admin regisretion",
			method:     http.MethodPost,
			path:       "/auth/admin_register",
			statusCode: 201,
			requestBody: db.Admin{Name: testAdmin.Name,
				FirstName: testAdmin.FirstName,
				LastName:  testAdmin.LastName,
				Password:  testAdmin.Password},
		},
		{
			name:       "Admin with this name already exist",
			method:     http.MethodPost,
			path:       "/auth/admin_register",
			statusCode: 409,
			requestBody: db.Admin{Name: testAdmin.Name,
				FirstName: testAdmin.FirstName,
				LastName:  testAdmin.LastName,
				Password:  testAdmin.Password},
		},
	}

	// connect db
	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	handlers.CreateFirstAdmin()

	firstAdmin := db.Admin{Name: "admin", Password: "1234567890"}

	adminTokens := adminLogin(firstAdmin, testServer)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsedRequestBody, _ := json.Marshal(tc.requestBody)
			testRequest := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(parsedRequestBody))
			testRequest.Header.Add("Content-Type", "application/json")
			testRequest.Header.Add("accessToken", adminTokens.AccessToken)

			fmt.Println(tc.requestBody)

			response, _ := testServer.Test(testRequest)

			receivedStatusCode := response.StatusCode
			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}
		})
	}
	db.DeleteBy[db.Admin]("name", firstAdmin.Name)

}

func TestAdminLoginHandler(t *testing.T) {
	if err := godotenv.Load("../dev.env"); err != nil {
		log.Fatalf("Error loading dev.env file")
	}

	config := db.DatabaseConfig{
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
		Host:     os.Getenv("DATABASE_HOST"),
	}

	testAdmin := db.Admin{Name: "greatestAdmin13378",
		FirstName: "Alexander",
		LastName:  "Nevsky",
		Password:  "superpassword"}

	testServer := fiber.New()
	controllers.AuthController(testServer)

	// connect db
	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	handlers.CreateFirstAdmin()

	firstAdmin := db.Admin{Name: "admin", Password: "1234567890"}

	adminTokens := adminLogin(firstAdmin, testServer)

	_ = adminAuth(testAdmin, adminTokens.AccessToken, testServer)

	tests := []struct {
		name        string
		method      string
		path        string
		requestBody interface{}
		statusCode  int
	}{
		{
			name:        "Admin Login",
			method:      http.MethodPost,
			path:        "/auth/admin_login",
			statusCode:  200,
			requestBody: db.Admin{Name: testAdmin.Name, Password: testAdmin.Password},
		},
		{
			name:        "Bad login",
			method:      http.MethodPost,
			path:        "/auth/admin_login",
			statusCode:  400,
			requestBody: db.Admin{Name: "", Password: testAdmin.Password},
		},
		{
			name:        "Bad password",
			method:      http.MethodPost,
			path:        "/auth/admin_login",
			statusCode:  400,
			requestBody: db.Admin{Name: testAdmin.Name, Password: ""},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parsedRequestBody, _ := json.Marshal(tc.requestBody)
			testRequest := httptest.NewRequest(tc.method, tc.path, bytes.NewReader(parsedRequestBody))
			testRequest.Header.Add("Content-Type", "application/json")

			fmt.Println(tc.requestBody)

			response, _ := testServer.Test(testRequest)

			receivedStatusCode := response.StatusCode
			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}

		})
	}
	db.DeleteBy[db.Admin]("name", firstAdmin.Name)
}

func adminLogin(admin db.Admin, testServer *fiber.App) utility.JWTTokenPair {
	loginBody, _ := json.Marshal(admin)
	loginRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8081/auth/admin_login", bytes.NewReader(loginBody))
	loginRequest.Header.Add("Content-Type", "application/json")

	response, err := testServer.Test(loginRequest)
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var adminTokens utility.JWTTokenPair
	_ = json.Unmarshal(body, &adminTokens)
	return adminTokens
}

func adminAuth(adminToRegister db.Admin, token string, testServer *fiber.App) db.Admin {
	authBody, _ := json.Marshal(adminToRegister)
	authRequest := httptest.NewRequest(http.MethodPost, "http://localhost:8081/auth/admin_register", bytes.NewReader(authBody))
	authRequest.Header.Add("Content-Type", "application/json")
	authRequest.Header.Add("accessToken", token)

	response, err := testServer.Test(authRequest)
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var authedAdmin db.Admin
	_ = json.Unmarshal(body, &authedAdmin)
	return authedAdmin
}
