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

type CreateGrowthRequest struct {
	HeightCm         float64 `json:"height_cm"`
	GrowthRateStatus string  `json:"growth_rate_status"`
}

// GetGrowthMetrics godoc
// @Summary      Get growth metrics for a plant
// @Description  Get all growth records for a specific plant
// @Tags         growth
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Plant ID"
// @Success      200 {array} models.GrowthMetric
// @Failure      401 {string} string "unauthorized"
// @Router       /plants/{id}/growth [get]
func GetGrowthMetrics(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		http.Error(w, "plant not found", http.StatusNotFound)
		return
	}

	var metrics []models.GrowthMetric
	config.Db.Where("plant_id = ?", plant.ID).Order("recorded_date desc").Find(&metrics)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// CreateGrowthMetric godoc
// @Summary      Record a growth metric
// @Description  Record a new growth measurement for a plant
// @Tags         growth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Plant ID"
// @Param        body body CreateGrowthRequest true "Growth payload"
// @Success      201 {object} models.GrowthMetric
// @Failure      400 {string} string "invalid request"
// @Failure      401 {string} string "unauthorized"
// @Router       /plants/{id}/growth [post]
func CreateGrowthMetric(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		http.Error(w, "plant not found", http.StatusNotFound)
		return
	}

	var req CreateGrowthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	metric := models.GrowthMetric{
		PlantID:          plant.ID,
		RecordedDate:     time.Now(),
		HeightCm:         req.HeightCm,
		GrowthRateStatus: models.GrowthRateStatus(req.GrowthRateStatus),
	}

	if result := config.Db.Create(&metric); result.Error != nil {
		http.Error(w, "failed to create growth metric", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(metric)
}
