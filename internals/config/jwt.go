package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecret string

func GetJwtSecret() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to load env file")
		return
	}

	JwtSecret = os.Getenv("JWT_SECRET_KEY")
}
