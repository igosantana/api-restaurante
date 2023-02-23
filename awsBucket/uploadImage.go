package awsBucket

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadImage(file *multipart.FileHeader) (string, error) {
	f, openErr := file.Open()
	defer f.Close()
	if openErr != nil {
		return "", openErr
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return "", err
	}
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	result, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("igo-rest-images"),
		Key:    aws.String(file.Filename),
		Body:   f,
		ACL:    "public-read",
	})
	if uploadErr != nil {
		return "", uploadErr
	}

	return result.Location, err
}
