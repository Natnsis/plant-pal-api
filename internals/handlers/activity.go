package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
	"plantPal/internals/models"
	"plantPal/internals/response"

	"github.com/gorilla/mux"
)

type CreateActivityRequest struct {
	ActivityType string `json:"activity_type"`
	Notes        string `json:"notes"`
	PhotoUrl     string `json:"photo_url"`
}

func GetActivities(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		response.Error(w, http.StatusNotFound, "plant not found")
		return
	}

	var activities []models.ActivityLog
	config.Db.Where("plant_id = ?", plant.ID).Order("logged_date desc").Find(&activities)

	response.JSON(w, http.StatusOK, activities)
}

func CreateActivity(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		response.Error(w, http.StatusNotFound, "plant not found")
		return
	}

	var req CreateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ActivityType == "" {
		response.Error(w, http.StatusBadRequest, "activity_type is required")
		return
	}

	activity := models.ActivityLog{
		PlantID:      plant.ID,
		ActivityType: models.ActivityType(req.ActivityType),
		LoggedDate:   time.Now(),
		Notes:        req.Notes,
		PhotoUrl:     req.PhotoUrl,
	}

	if result := config.Db.Create(&activity); result.Error != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create activity")
		return
	}

	response.JSON(w, http.StatusCreated, activity)
}
