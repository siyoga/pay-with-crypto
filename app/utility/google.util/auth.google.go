package utility

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/grokify/go-pkce"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//"https://www.googleapis.com/auth/userinfo.email"

func ConfigGoogle() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"openid", "email", "https://www.googleapis.com/auth/userinfo.profile"}, // you can use other scopes to get more data
		Endpoint: google.Endpoint,
	}
	return conf
}

func GetTokens(token string, PCKECode string) (*GoogleOauthToken, error) {
	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", token)
	values.Add("client_id", os.Getenv("GOOGLE_CLIENT"))
	values.Add("client_secret", os.Getenv("GOOGLE_SECRET"))
	values.Add("redirect_uri", os.Getenv("GOOGLE_REDIRECT_URL"))
	values.Add("code_verifier", PCKECode)

	query := values.Encode()

	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", bytes.NewBufferString(query))
	if err != nil {
		panic(err)

	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve token")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleOauthTokenRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleOauthTokenRes); err != nil {
		return nil, err
	}

	tokenBody := &GoogleOauthToken{
		Access_token: GoogleOauthTokenRes["access_token"].(string),
		Id_token:     GoogleOauthTokenRes["id_token"].(string),
	}

	return tokenBody, nil
}

func GetUserData(tokens *GoogleOauthToken) (*GoogleUserResult, error) {
	rootUrl := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", tokens.Access_token)

	req, err := http.NewRequest("GET", rootUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.Id_token))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GoogleUserRes map[string]interface{}
	if err := json.Unmarshal(resBody.Bytes(), &GoogleUserRes); err != nil {

		return nil, err
	}

	userBody := &GoogleUserResult{
		Id:             GoogleUserRes["id"].(string),
		Email:          GoogleUserRes["email"].(string),
		Name:           GoogleUserRes["name"].(string),
		Verified_email: GoogleUserRes["verified_email"].(bool),
		Given_name:     GoogleUserRes["given_name"].(string),
		Picture:        GoogleUserRes["picture"].(string),
		Locale:         GoogleUserRes["locale"].(string),
	}

	return userBody, nil
}

func CreatePKCE() *PKCE {
	NewPKCE := PKCE{
		pkce.NewCodeVerifier(),
		""}

	NewPKCE.CodeChallenge = (pkce.CodeChallengeS256(NewPKCE.CodeVerifier))

	ThePKCE = &NewPKCE
	return ThePKCE
}
