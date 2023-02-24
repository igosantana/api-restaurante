package awsBucket

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"

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
	key := file.Filename + time.Now().String()
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	result, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("restaurant-images-igo"),
		Key:    aws.String(key),
		Body:   f,
		ACL:    "public-read",
	})
	fmt.Println("uploadErro", uploadErr)
	if uploadErr != nil {
		return "", uploadErr
	}
	return result.Location, err
}
