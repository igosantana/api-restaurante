package helper

import (
	"context"
	"mime/multipart"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func ImageUploadHelper(file *multipart.FileHeader) (string, error) {
	f, openErr := file.Open()
	defer f.Close()
	if openErr != nil {
		return "", openErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//create cloudnary instance
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	cloudAPIKey := os.Getenv("CLOUDINARY_API_KEY")
	cloudAPISecret := os.Getenv("CLOUDINARY_API_SECRET")
	cld, err := cloudinary.NewFromParams(cloudName, cloudAPIKey, cloudAPISecret)
	if err != nil {
		return "", err
	}

	//upload file
	uploadFolder := os.Getenv("CLOUDINARY_UPLOAD_FOLDER")
	uploadParam, err := cld.Upload.Upload(ctx, f, uploader.UploadParams{Folder: uploadFolder})
	if err != nil {
		return "", err
	}

	return uploadParam.SecureURL, nil
}
