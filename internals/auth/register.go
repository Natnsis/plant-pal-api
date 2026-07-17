package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"plantPal/internals/config"
	"plantPal/internals/models"
	"plantPal/internals/response"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with full name, email, password, and phone number
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "User registration payload"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  response.ErrorResponse
// @Failure      409   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FullName == "" || req.Email == "" || req.Password == "" || req.PhoneNumber == "" {
		response.Error(w, http.StatusBadRequest, "full_name, email, password, and phone_number are required")
		return
	}

	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		response.Error(w, http.StatusBadRequest, "invalid email format")
		return
	}

	if len(req.Password) < 8 {
		response.Error(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user := models.User{
		Email:       req.Email,
		Password:    string(hashedPassword),
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}

	if result := config.Db.Create(&user); result.Error != nil {
		response.Error(w, http.StatusConflict, "email or phone number already exists")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"id":         user.ID,
		"full_name":  user.FullName,
		"email":      user.Email,
		"phone":      user.PhoneNumber,
	})
}
