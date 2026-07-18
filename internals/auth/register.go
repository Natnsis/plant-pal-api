package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"plantPal/internals/config"
	"plantPal/internals/models"
	"plantPal/internals/response"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with full name, email, password, and phone number
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "User registration payload"
// @Success      201   {object}  map[string]string
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

	if req.FullName == "" || req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "full_name, email, and password are required")
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
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
	}

	if result := config.Db.Create(&user); result.Error != nil {
		response.Error(w, http.StatusConflict, "email or phone number already exists")
		return
	}

	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.FullName,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.FullName,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 15).Unix(),
	}

	unsignedAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	unsignedRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessToken, err := unsignedAccess.SignedString([]byte(config.JwtSecret))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	refreshToken, err := unsignedRefresh.SignedString([]byte(config.JwtSecret))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to generate refresh token")
		return
	}

	dbToken := models.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 15),
	}
	if result := config.Db.Create(&dbToken); result.Error != nil {
		response.Error(w, http.StatusInternalServerError, "failed to store refresh token")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
