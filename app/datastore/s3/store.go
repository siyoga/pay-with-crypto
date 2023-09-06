package s3

import (
	"pay-with-crypto/app/utility"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func S3Connect() (*session.Session, bool) {
	endpoint := utility.GetEnv("S3_HOST", "")
	accessKey := utility.GetEnv("S3_ACCESSKEY", "")
	secretKey := utility.GetEnv("S3_SECRETKEY", "")
	region := utility.GetEnv("S3_REGION", "ru-1")

	sess, err := session.NewSession(
		&aws.Config{
			Endpoint: aws.String(endpoint),
			Region:   aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKey,
				secretKey,
				"",
			),
		})

	if err != nil {
		return nil, false
	}

	return sess, true
}
