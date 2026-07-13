package auth

import (
	"encoding/json"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/models"
)

// request body type
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	// decode and assign request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and Password required", http.StatusBadRequest)
		return
	}

	// fetch user with that specific email
	var user models.User
	result := config.Db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	// bcrypt password comparision

	// generate tokens

	// send back to user
}
