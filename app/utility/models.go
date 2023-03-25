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

	Status struct {
		ID     uuid.UUID `json:"id" gorm:"type:uuid"`
		Status bool      `json:"status" gorm:"bool"`
	}

	Message struct {
		Text string `json:"text"`
	}
)
