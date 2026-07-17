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

// GetGrowthMetrics godoc
// @Summary      Get growth metrics
// @Description  Retrieve all growth metrics for a specific plant
// @Tags         growth
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Plant ID"
// @Success      200  {array}   models.GrowthMetric
// @Failure      404  {object}  response.ErrorResponse
// @Router       /plants/{id}/growth [get]
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

// CreateGrowthMetric godoc
// @Summary      Create a growth metric
// @Description  Record a new growth measurement for a plant
// @Tags         growth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string              true  "Plant ID"
// @Param        body  body      CreateGrowthRequest  true  "Growth measurement"
// @Success      201   {object}  models.GrowthMetric
// @Failure      400   {object}  response.ErrorResponse
// @Failure      404   {object}  response.ErrorResponse
// @Failure      500   {object}  response.ErrorResponse
// @Router       /plants/{id}/growth [post]
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
