package config

import "os"

var GeminiAPIKey string

func GetGeminiAPIKey() {
	GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
}
