package utility

import "github.com/gofrs/uuid"

type (
	JWTTokenPair struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	Bucket struct {
		BucketName string
		Location   string
	}

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
	
  Status struct {
		ID     uuid.UUID `json:"id" gorm:"type:uuid"`
		Status bool      `json:"status" gorm:"bool"`
	}
)
