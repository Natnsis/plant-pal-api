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

type CreateGrowthRequest struct {
	HeightCm         float64 `json:"height_cm"`
	GrowthRateStatus string  `json:"growth_rate_status"`
}

func GetGrowthMetrics(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		response.Error(w, http.StatusNotFound, "plant not found")
		return
	}

	var metrics []models.GrowthMetric
	config.Db.Where("plant_id = ?", plant.ID).Order("recorded_date desc").Find(&metrics)

	response.JSON(w, http.StatusOK, metrics)
}

func CreateGrowthMetric(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		response.Error(w, http.StatusNotFound, "plant not found")
		return
	}

	var req CreateGrowthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	metric := models.GrowthMetric{
		PlantID:          plant.ID,
		RecordedDate:     time.Now(),
		HeightCm:         req.HeightCm,
		GrowthRateStatus: models.GrowthRateStatus(req.GrowthRateStatus),
	}

	if result := config.Db.Create(&metric); result.Error != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create growth metric")
		return
	}

	response.JSON(w, http.StatusCreated, metric)
}
