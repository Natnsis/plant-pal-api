package config

import "os"

var GoogleClientID string

func GetGoogleConfig() {
	GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
}
