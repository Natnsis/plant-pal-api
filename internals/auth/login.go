package auth

import (
	"encoding/json"
	"net/http"
)

// request body type
type RequestType struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	// check method
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusBadGateway)
		return
	}

	// get user data
	var req RequestType
	json.NewDecoder(r.Body).Decode(&req)

	// validating email and password
	if req.Email == "" || req.Password == "" {
		http.Error(w, "some credentials are missing", http.StatusBadGateway)
		return
	}

	// compare the hashed password

	// generate tokens

	// send back to user
}
