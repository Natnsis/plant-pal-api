package auth

import (
	"encoding/json"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/models"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func AddUser() {
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusBadGateway)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "request body is empty", http.StatusNoContent)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email or password required", http.StatusNoContent)
		return
	}

	var exists bool
	var user User
	if err := config.Db.Find(&user, req.Email == user.Email); err != nil {
		exists = true
	} else {
		exists = false
	}

	if !exists {
		http.Error(w, "email already exits", http.StatusConflict)
		return
	}

	// hash password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "unable to hash password", http.StatusInternalServerError)
		return
	}

	userData := models.User{
		FirstName:   "",
		LastName:    "",
		PhoneNumber: "",
		Email:       req.Email,
		Password:    string(hashedPassword),
	}

	result := config.Db.Create(&userData)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
}
