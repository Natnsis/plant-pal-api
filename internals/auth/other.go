package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"plantPal/internals/config"
	"plantPal/internals/models"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Refresh godoc
// @Summary      Refresh tokens
// @Description  Exchange a refresh token for a new access and refresh token pair (single-use rotation)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RefreshRequest  true  "Refresh token payload"
// @Success      200   {object}  map[string]string
// @Failure      400   {string}  string  "invalid request body or missing refresh_token"
// @Failure      401   {string}  string  "invalid, expired, or revoked refresh token"
// @Router       /refresh [post]
func Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}

	// Parse and validate the token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(config.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// Check if token exists in database and is not revoked
	var dbToken models.RefreshToken
	if result := config.Db.Where("token = ? AND revoked = false", req.RefreshToken).First(&dbToken); result.Error != nil {
		http.Error(w, "refresh token not found or revoked", http.StatusUnauthorized)
		return
	}

	// Check expiry
	if time.Now().After(dbToken.ExpiresAt) {
		config.Db.Delete(&dbToken)
		http.Error(w, "refresh token expired", http.StatusUnauthorized)
		return
	}

	// Revoke the old token (single-use rotation)
	config.Db.Model(&dbToken).Update("revoked", true)

	// Generate new token pair
	claims, _ := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)

	accessClaims := jwt.MapClaims{
		"user_id": uint(userID),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id": uint(userID),
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

	// Store new refresh token
	newDbToken := models.RefreshToken{
		Token:     refreshToken,
		UserID:    uint(userID),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 15),
	}
	if result := config.Db.Create(&newDbToken); result.Error != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Logout godoc
// @Summary      Logout a user
// @Description  Revoke a refresh token to invalidate a session
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      LogoutRequest  true  "Logout payload"
// @Success      200   {object}  map[string]string
// @Failure      400   {string}  string  "invalid request body or missing refresh_token"
// @Failure      404   {string}  string  "refresh token not found or already revoked"
// @Router       /logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}

	// Revoke the refresh token
	result := config.Db.Model(&models.RefreshToken{}).
		Where("token = ? AND revoked = false", req.RefreshToken).
		Update("revoked", true)

	if result.RowsAffected == 0 {
		http.Error(w, "refresh token not found or already revoked", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out successfully"})
}
