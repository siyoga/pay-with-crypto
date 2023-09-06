package s3

import (
	"mime/multipart"
	"pay-with-crypto/app/utility"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadFile(file multipart.File, fileName string) (*string, bool) {
	bucket := utility.GetEnv("S3_BUCKET", "")
	sess, isOk := S3Connect()
	uploader := s3manager.NewUploader(sess)

	if !isOk {
		utility.Error(nil, "Failed to connect to the S3 Bucket")
		return nil, false
	}

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		utility.Error(err, "Failed to upload image to the S3 Bucket")
		return nil, false
	}

	fileLink := "https://s3.pay-with-crypto.xyz/" + fileName

	return &fileLink, true
}

func DeleteImage(fileName string) bool {
	bucket := utility.GetEnv("S3_BUCKET", "")
	sess, isOk := S3Connect()

	if !isOk {
		utility.Error(nil, "Failed to connect to the S3 Bucket")
		return false
	}

	svc := s3.New(sess)
	objectInfo := s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}

	_, err := svc.DeleteObject(&objectInfo)

	if err != nil {
		utility.Error(err, "Failed to delete image from the S3 Bucket")
		return false
	}

	return true
}
