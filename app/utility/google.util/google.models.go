package utility

var (
	ThePKCE *PKCE
)

type (
	GoogleOauthToken struct {
		Access_token string
		Id_token     string
	}

	GoogleUserResult struct {
		Id             string
		Email          string
		Verified_email bool
		Name           string
		Given_name     string
		Family_name    string
		Picture        string
		Locale         string
	}

	PKCE struct {
		CodeVerifier  string
		CodeChallenge string
	}

	GoogleIDToken struct {
		Iss           string `json:"iss"`
		Sub           string `json:"sub"`
		Azp           string `json:"azp"`
		Aud           string `json:"aud"`
		Iat           string `json:"iat"`
		Exp           string `json:"exp"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Locale        string `json:"locale"`
	}

	GoogleErrorResponse struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
)
