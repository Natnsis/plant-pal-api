package config

import "os"

var JwtSecret string

func GetJwtSecret() {
	JwtSecret = os.Getenv("JWT_SECRET_KEY")
}
