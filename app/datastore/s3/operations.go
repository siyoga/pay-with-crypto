package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"pay-with-crypto/app/utility"

	"github.com/minio/minio-go/v7"
)

func UploadFile(file *multipart.FileHeader, buffer multipart.File, bucketName string, fileName string) (*string, bool) {
	ctx := context.Background()
	minioClient, isOk := MinioConnect()

	if !isOk {
		return nil, false
	}

	ext := filepath.Ext(file.Filename)
	fileName = fileName + ext
	fileBuffer := buffer
	contentType := file.Header["Content-Type"][0]
	fileSize := file.Size

	_, err := minioClient.PutObject(ctx, bucketName, fileName, fileBuffer, fileSize, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		fmt.Println(err)
		utility.Error(err, "Upload File")
		return nil, false
	}

	return &fileName, true
}
