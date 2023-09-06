package utility

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt"
)

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

	RegisterInfoRequest struct {
		ViaGoogle bool   `json:"viaGoogle"`
		Email     string `json:"email"`
	}

	UpdateInfoRequest struct {
		Link  string `json:"link"`
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	Claims struct {
		Sub       string    `json:"sub"`
		Generated time.Time `json:"generated"`
		jwt.StandardClaims
	}
)
