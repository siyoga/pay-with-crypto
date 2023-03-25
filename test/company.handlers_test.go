package test

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
	"github.com/joho/godotenv"
)

func TestSearchCompanyById(t *testing.T) {
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

	testCompany := db.Company{
		Name:          "TestSearchCompanyById",
		Password:      "TestSearchCompanyById_PASSWORD",
		Mail:          "testtt@kekw.com",
		LinkToCompany: "http://linktt.com"}

	company := auth(testCompany, testServer)

	tests := []struct {
		name       string
		method     string
		path       string
		statusCode int
	}{
		{
			name:       "Successful pipeline",
			method:     http.MethodGet,
			path:       "http://127.0.0.1:8081/company/search/userid?id=" + company.ID.String(),
			statusCode: 200,
		},
		{
			name:       "No id pipeline",
			method:     http.MethodGet,
			path:       "http://127.0.0.1:8081/company/search/userid?id=",
			statusCode: 400,
		},
		{
			name:       "No user bound to id pipeline",
			method:     http.MethodGet,
			path:       "http://127.0.0.1:8081/company/search/userid?id=" + "df2b73fs-sdse-4cee-a13e-7e60766f7992",
			statusCode: 404,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			testRequest := httptest.NewRequest(tc.method, tc.path, nil)

			response, _ := testServer.Test(testRequest)
			receivedStatusCode := response.StatusCode
			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}

		})
	}
}

func TestCompanyUploadPicture(t *testing.T) {
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

	testCompany := db.Company{
		Name:          "TEST_COMPANY_NAME_Upload_Picture",
		Password:      "TEST_COMPANY_NAME_Upload_Picture_PASSWORD",
		Mail:          "testtt@kekw.com",
		LinkToCompany: "http://linktt.com"}

	company := auth(testCompany, testServer)

	tokens := login(db.Company{Name: testCompany.Name, Password: testCompany.Password}, testServer)

	fullPath, _ := filepath.Abs("../test/logos/logo_company.jpg")

	tests := []struct {
		name        string
		method      string
		path        string
		requestBody db.Company
		statusCode  int
		pathImage   string
	}{
		{
			name:        "Successful pipeline",
			method:      http.MethodPost,
			path:        "http://127.0.0.1:8081/company/uploadLogo",
			statusCode:  204,
			requestBody: db.Company{ID: company.ID, Image: ""},
			pathImage:   fullPath,
		},
		{
			name:        "No image pipeline",
			method:      http.MethodPost,
			path:        "http://127.0.0.1:8081/company/uploadLogo",
			statusCode:  400,
			requestBody: db.Company{ID: company.ID, Image: ""},
			pathImage:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bw, buf, err := createMultiform(tc.pathImage, company.ID.String())

			if err != nil {
				fmt.Print(err) //Not sure that how it should be
			}
			testRequest := httptest.NewRequest(tc.method, tc.path, buf)
			testRequest.Header.Add("Content-Type", bw.FormDataContentType())
			testRequest.Header.Add("accessToken", tokens.AccessToken)

			response, err := testServer.Test(testRequest)
			if err != nil {
				fmt.Println(err)
			}

			receivedStatusCode := response.StatusCode

			if receivedStatusCode != tc.statusCode {
				t.Errorf("StatusCode was incorrect, got: %d, want: %d", receivedStatusCode, tc.statusCode)
			}
		})
	}
}

func auth(testCompany db.Company, testServer *fiber.App) db.Company {
	authBody, _ := json.Marshal(testCompany)
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
	return company
}

func login(testCompany db.Company, testServer *fiber.App) utility.JWTTokenPair {
	loginBody, _ := json.Marshal(testCompany)
	loginRequest := httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8081/auth/login", bytes.NewReader(loginBody))
	loginRequest.Header.Add("Content-Type", "application/json")
	response, err := testServer.Test(loginRequest)

	if err != nil {
		fmt.Println(err)
	}
	body, _ := io.ReadAll(response.Body)

	var tokens utility.JWTTokenPair
	_ = json.Unmarshal(body, &tokens)
	return tokens
}

func createMultiform(path string, ID string) (*multipart.Writer, *bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)

	p1w, _ := bw.CreateFormField("id")
	p1w.Write([]byte(ID))

	f, err := os.Open(path)
	if err != nil {
		f.Close()
		bw.Close()
		return bw, buf, err
	}
	defer f.Close()

	_, fileName := filepath.Split(path)
	fw1, _ := bw.CreateFormFile("companyLogo", fileName)
	io.Copy(fw1, f)
	bw.Close()

	return bw, buf, nil
}
