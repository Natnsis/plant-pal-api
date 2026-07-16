package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"plantPal/internals/config"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	result, err := config.Cloudinary.Upload.Upload(
		ctx,
		string(fileBytes),
		uploader.UploadParams{
			Folder: "plantpal",
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload to cloudinary: %w", err)
	}

	return result.SecureURL, nil
}
