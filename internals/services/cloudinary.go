package services

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"

	"plantPal/internals/config"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(file multipart.File, _ *multipart.FileHeader) (string, error) {
	ctx := context.Background()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	result, err := config.Cloudinary.Upload.Upload(
		ctx,
		bytes.NewReader(buf.Bytes()),
		uploader.UploadParams{
			Folder: "plantpal",
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload to cloudinary: %w", err)
	}

	return result.SecureURL, nil
}
