package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var Cloudinary *cloudinary.Cloudinary

func GetCloudinaryConfig() {
	cldURL := os.Getenv("CLOUDINARY_URL")
	if cldURL == "" {
		log.Fatal("no CLOUDINARY_URL found")
	}

	cld, err := cloudinary.NewFromURL(cldURL)
	if err != nil {
		log.Fatal("failed to init cloudinary: ", err)
	}

	Cloudinary = cld
}
