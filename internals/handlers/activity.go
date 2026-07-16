package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
	"plantPal/internals/models"

	"github.com/gorilla/mux"
)

type CreateActivityRequest struct {
	ActivityType string `json:"activity_type"`
	Notes        string `json:"notes"`
	PhotoUrl     string `json:"photo_url"`
}

// GetActivities godoc
// @Summary      Get activity log for a plant
// @Description  Get all activities logged for a specific plant
// @Tags         activity
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Plant ID"
// @Success      200 {array} models.ActivityLog
// @Failure      401 {string} string "unauthorized"
// @Router       /plants/{id}/activities [get]
func GetActivities(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		http.Error(w, "plant not found", http.StatusNotFound)
		return
	}

	var activities []models.ActivityLog
	config.Db.Where("plant_id = ?", plant.ID).Order("logged_date desc").Find(&activities)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}

// CreateActivity godoc
// @Summary      Log an activity
// @Description  Log a care activity for a plant
// @Tags         activity
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Plant ID"
// @Param        body body CreateActivityRequest true "Activity payload"
// @Success      201 {object} models.ActivityLog
// @Failure      400 {string} string "invalid request"
// @Failure      401 {string} string "unauthorized"
// @Router       /plants/{id}/activities [post]
func CreateActivity(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		http.Error(w, "plant not found", http.StatusNotFound)
		return
	}

	var req CreateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ActivityType == "" {
		http.Error(w, "activity_type is required", http.StatusBadRequest)
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
		http.Error(w, "failed to create activity", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(activity)
}


