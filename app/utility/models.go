package utility

type (
	JWTTokenPair struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	Bucket struct {
		BucketName string
		Location   string
	}
)
