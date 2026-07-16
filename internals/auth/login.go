package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"plantPal/internals/config"
	"plantPal/internals/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login godoc
// @Summary      Login a user
// @Description  Authenticate with email and password to receive access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "Login credentials"
// @Success      200   {object}  map[string]string
// @Failure      400   {string}  string  "invalid request body or missing fields"
// @Failure      401   {string}  string  "invalid credentials"
// @Router       /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	var user models.User
	if result := config.Db.Where("email = ?", req.Email).First(&user); result.Error != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 15).Unix(),
	}

	unsignedAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	unsignedRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessToken, err := unsignedAccess.SignedString([]byte(config.JwtSecret))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	refreshToken, err := unsignedRefresh.SignedString([]byte(config.JwtSecret))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Store refresh token in database
	dbToken := models.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 15),
	}
	if result := config.Db.Create(&dbToken); result.Error != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
