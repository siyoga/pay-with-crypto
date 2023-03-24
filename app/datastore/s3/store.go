package s3

import (
	"context"
	"fmt"
	"os"
	"pay-with-crypto/app/utility"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func MinioConnect() (*minio.Client, bool) {
	ctx := context.Background()
	endpoint := os.Getenv("S3_HOST")
	accessKeyId := os.Getenv("S3_USER")
	secretAccessKey := os.Getenv("S3_PASSWORD")

	fmt.Println("s3 endpoint: " + endpoint)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		utility.Error(err, "Minio Init")
		return nil, false
	}

	var buckets [3]utility.Bucket

	buckets[0] = utility.Bucket{BucketName: os.Getenv("S3_BUCKET_CARD_LOGO"), Location: "us-east-1"}
	buckets[1] = utility.Bucket{BucketName: os.Getenv("S3_BUCKET_CARD_IMAGES"), Location: "us-east-1"}
	buckets[2] = utility.Bucket{BucketName: os.Getenv("S3_BUCKET_COMPANY_LOGOS"), Location: "us-east-1"}

	for i := 0; i < 3; i++ {
		err := minioClient.MakeBucket(ctx, buckets[i].BucketName, minio.MakeBucketOptions{Region: buckets[i].Location})

		if err != nil {
			fmt.Println(err)
			exist, errBucketExists := minioClient.BucketExists(ctx, buckets[i].BucketName)
			if errBucketExists == nil && exist {
				fmt.Printf("Bucket %s already created. \n", buckets[i].BucketName)
			} else {
				fmt.Printf("Cannot create bucket %s \n", buckets[i].BucketName)
				utility.Error(err, "Make Bucket")
			}
		} else {
			fmt.Printf("Bucket %s successfully created. \n", buckets[i].BucketName)
		}
	}

	return minioClient, true
}
