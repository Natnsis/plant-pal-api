package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"plantPal/internals/config"
	"plantPal/internals/models"
	"plantPal/internals/response"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)

type GoogleLoginRequest struct {
	IDToken string `json:"id_token"`
}

type GoogleLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID                   uint   `json:"id"`
	FullName             string `json:"full_name"`
	Email                string `json:"email"`
	CareStreakDays       int    `json:"care_streak_days"`
	TotalTaskDone        int    `json:"total_task_done"`
	TotalJournalInjuries int    `json:"total_journial_injuries"`
}

// LoginWithGoogle godoc
// @Summary      Login or register with Google
// @Description  Authenticate using a Google ID token. Creates a new user if the email doesn't exist.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      GoogleLoginRequest  true  "Google ID token"
// @Success      200   {object}  GoogleLoginResponse
// @Failure      400   {object}  response.ErrorResponse
// @Failure      401   {object}  response.ErrorResponse
// @Failure      409   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /auth/google [post]
func LoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	var req GoogleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.IDToken == "" {
		response.Error(w, http.StatusBadRequest, "id_token is required")
		return
	}

	payload, err := idtoken.Validate(r.Context(), req.IDToken, config.GoogleClientID)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid Google ID token")
		return
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	if email == "" {
		response.Error(w, http.StatusBadRequest, "email not provided by Google")
		return
	}

	var user models.User
	result := config.Db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("google_oauth_"+email), bcrypt.DefaultCost)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to create user")
			return
		}

		user = models.User{
			Email:    email,
			Password: string(hashedPassword),
			FullName: name,
		}

		if result := config.Db.Create(&user); result.Error != nil {
			response.Error(w, http.StatusConflict, "failed to create user")
			return
		}
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

	response.JSON(w, http.StatusOK, GoogleLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserResponse{
			ID:                   user.ID,
			FullName:             user.FullName,
			Email:                user.Email,
			CareStreakDays:       user.CareStreakDays,
			TotalTaskDone:        user.TotalTaskDone,
			TotalJournalInjuries: user.TotalJournalInjuries,
		},
	})
}
