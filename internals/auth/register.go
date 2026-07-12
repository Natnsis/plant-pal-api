package auth

import (
	"encoding/json"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/models"

	"golang.org/x/crypto/bcrypt"
)

// list the request data types
type RequestData struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	// check method
	if r.Method != http.MethodPost {
		http.Error(w, "wrong request method used", http.StatusMethodNotAllowed)
		return
	}

	// get userdata
	var req RequestData
	json.NewDecoder(r.Body).Decode(&req)

	// validate user data
	if req.FullName == "" || req.Email == "" || req.Password == "" || req.PhoneNumber == "" {
		http.Error(w, "missing parameter", http.StatusBadRequest)
		return
	}

	// hash user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "unable to hash the password", http.StatusBadRequest)
		return
	}

	// put user data into table with hashed data
	UserData := models.User{
		Email:       req.Email,
		Password:    string(hashedPassword),
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}

	reslut := config.Db.Create(&UserData)
	if reslut.Error != nil {
		http.Error(w, "unable to hash the password", http.StatusBadRequest)
		return
	}

	// return value
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UserData)
}
