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
)
