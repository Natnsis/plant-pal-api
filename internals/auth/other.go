package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"plantPal/internals/config"
	"plantPal/internals/models"
	"plantPal/internals/response"

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
// @Failure      400   {object}  response.ErrorResponse
// @Failure      401   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /refresh [post]
func Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(config.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		response.Error(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	var dbToken models.RefreshToken
	if result := config.Db.Where("token = ? AND revoked = false", req.RefreshToken).First(&dbToken); result.Error != nil {
		response.Error(w, http.StatusUnauthorized, "refresh token not found or revoked")
		return
	}

	if time.Now().After(dbToken.ExpiresAt) {
		config.Db.Delete(&dbToken)
		response.Error(w, http.StatusUnauthorized, "refresh token expired")
		return
	}

	config.Db.Model(&dbToken).Update("revoked", true)

	claims, _ := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)
	name, _ := claims["name"].(string)
	email, _ := claims["email"].(string)

	accessClaims := jwt.MapClaims{
		"user_id": uint(userID),
		"name":    name,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id": uint(userID),
		"name":    name,
		"email":   email,
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

	newDbToken := models.RefreshToken{
		Token:     refreshToken,
		UserID:    uint(userID),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 15),
	}
	if result := config.Db.Create(&newDbToken); result.Error != nil {
		response.Error(w, http.StatusInternalServerError, "failed to store refresh token")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
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
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Router       /logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	result := config.Db.Model(&models.RefreshToken{}).
		Where("token = ? AND revoked = false", req.RefreshToken).
		Update("revoked", true)

	if result.RowsAffected == 0 {
		response.Error(w, http.StatusNotFound, "refresh token not found or already revoked")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}
