package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"pay-with-crypto/app/controllers"
	db "pay-with-crypto/app/datastore"
	"pay-with-crypto/app/utility"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/joho/godotenv"
)

func TestSearchCompanyById(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		statusCode int
	}{
		{
			name:       "Successful pipeline",
			method:     http.MethodGet,
			path:       "http://127.0.0.1:8081/company/search/userid?id=",
			statusCode: 200,
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
	controllers.CompanyController(testServer)

	// connect db
	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	authBody, _ := json.Marshal(
		db.Company{Name: "TestSearchCompanyById", Password: "TestSearchCompanyById_PASSWORD",
			Mail: "testtt@kekw.com", LinkToCompany: "http://linktt.com"})
	authRequest := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8081/auth/register", bytes.NewReader(authBody))
	authRequest.Header.Add("Content-Type", "application/json")

	response, err := testServer.Test(authRequest)
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var company db.Company
	_ = json.Unmarshal(body, &company)
	fmt.Println(company)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.path += company.ID.String()
			fmt.Println(tc.path)
			testRequest := httptest.NewRequest(tc.method, tc.path, nil)

			response, _ := testServer.Test(testRequest)
			fmt.Println(response.Status)
			receivedStatusCode := response.StatusCode
			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}

			body, err := io.ReadAll(response.Body)

			if err != nil {
				t.Error("Incorrect response body, check the sent json")
			}

			var createdTestUser db.Company

			if err := json.Unmarshal(body, &createdTestUser); err != nil {
				fmt.Println(createdTestUser)
				t.Errorf("Can't convert to Company struct. Error description: %s", err.Error())
			}

		})
	}
}

func TestCompanyUploadPicture(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		requestBody db.Company
		statusCode  int
	}{
		{
			name:        "Successful pipeline",
			method:      http.MethodPost,
			path:        "http://127.0.0.1:8081/company/uploadLogo",
			statusCode:  204,
			requestBody: db.Company{ID: uuid.Nil, Image: ""},
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
	controllers.CompanyController(testServer)

	// connect db
	db.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s", config.Host, config.User, config.Password, config.Database))

	authBody, _ := json.Marshal(
		db.Company{Name: "TEST_COMPANY_NAME_Upload_Picture", Password: "TEST_COMPANY_NAME_Upload_Picture_PASSWORD",
			Mail: "testtt@kekw.com", LinkToCompany: "http://linktt.com"})
	authRequest := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8081/auth/register", bytes.NewReader(authBody))
	authRequest.Header.Add("Content-Type", "application/json")

	response, err := testServer.Test(authRequest)
	if err != nil {
		fmt.Println(err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var company db.Company
	_ = json.Unmarshal(body, &company)

	loginBody, _ := json.Marshal(db.Company{Name: "TEST_COMPANY_NAME_Upload_Picture", Password: "TEST_COMPANY_NAME_Upload_Picture_PASSWORD"})
	loginRequest := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8081/auth/login", bytes.NewReader(loginBody))
	loginRequest.Header.Add("Content-Type", "application/json")
	response, _ = testServer.Test(loginRequest)

	if err != nil {
		fmt.Println(err)
	}
	body, err = io.ReadAll(response.Body)

	var tokens utility.JWTTokenPair
	_ = json.Unmarshal(body, &tokens)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.requestBody.ID = company.ID

			buf := new(bytes.Buffer)
			bw := multipart.NewWriter(buf)

			f, err := os.Open("/home/pyegor/Develop/pay-with-crypto/test/ABOBUS.png")
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()

			p1w, _ := bw.CreateFormField("id")
			p1w.Write([]byte(company.ID.String()))

			_, fileName := filepath.Split("/home/pyegor/Develop/pay-with-crypto/test/ABOBUS.png")
			fw1, _ := bw.CreateFormFile("companyLogo", fileName)
			io.Copy(fw1, f)
			bw.Close()

			// parsedRequestBody, _ := json.Marshal(tc.requestBody)
			testRequest := httptest.NewRequest(tc.method, tc.path, buf)
			testRequest.Header.Add("Content-Type", bw.FormDataContentType())
			testRequest.Header.Add("accessToken", tokens.AccessToken)

			response, _ = testServer.Test(testRequest)
			receivedStatusCode := response.StatusCode
			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}
		})
	}
}
