package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"plantPal/internals/config"
	"plantPal/internals/models"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.FullName == "" || req.Email == "" || req.Password == "" || req.PhoneNumber == "" {
		http.Error(w, "full_name, email, password, and phone_number are required", http.StatusBadRequest)
		return
	}

	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		http.Error(w, "password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email:       req.Email,
		Password:    string(hashedPassword),
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}

	if result := config.Db.Create(&user); result.Error != nil {
		http.Error(w, "email or phone number already exists", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        user.ID,
		"full_name": user.FullName,
		"email":     user.Email,
		"phone":     user.PhoneNumber,
	})
}
