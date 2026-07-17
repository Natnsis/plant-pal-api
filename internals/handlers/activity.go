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

// GetActivities godoc
// @Summary      Get activity log
// @Description  Retrieve all activity logs for a specific plant
// @Tags         activities
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Plant ID"
// @Success      200  {array}   models.ActivityLog
// @Failure      404  {object}  response.ErrorResponse
// @Router       /plants/{id}/activities [get]
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

// CreateActivity godoc
// @Summary      Log an activity
// @Description  Record a new activity for a plant (watered, fertilized, repotted, etc.)
// @Tags         activities
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                  true  "Plant ID"
// @Param        body  body      CreateActivityRequest  true  "Activity details"
// @Success      201   {object}  models.ActivityLog
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /plants/{id}/activities [post]
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
