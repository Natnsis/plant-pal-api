package auth

import (
	"net/http"

	"plantPal/internals/config"

	"github.com/golang-jwt/jwt/v5"
)

func Logout(w http.ResponseWriter, r *http.Request) {
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenString := r.Header.Get("Refresh-Token")
	if refreshTokenString == "" {
		http.Error(w, "refresh token required", http.StatusUnauthorized)
		return
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return config.JwtSecret, nil
	}

	token, err := jwt.Parse(refreshTokenString, keyFunc)
	if err != nil {
		http.Error(w, "refresh token required", http.StatusUnauthorized)
		return
	}
}
