package handlers

import (
	"encoding/json"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
	"plantPal/internals/models"

	"github.com/gorilla/mux"
)

type UpdateCarePlanRequest struct {
	WateringFrequencyDays *uint   `json:"watering_frequency_days"`
	WateringAmount        *string `json:"watering_amount"`
	WateringMethod        *string `json:"watering_method"`
	WateringTips          *string `json:"watering_tips"`
	LightRequirement      *string `json:"light_requirement"`
	HumidityRequirement   *string `json:"humidity_requirement"`
}

// GetCarePlan godoc
// @Summary      Get care plan for a plant
// @Description  Get the care plan associated with a plant
// @Tags         care-plan
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Plant ID"
// @Success      200 {object} models.CarePlan
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "care plan not found"
// @Router       /plants/{id}/care-plan [get]
func GetCarePlan(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		http.Error(w, "plant not found", http.StatusNotFound)
		return
	}

	var carePlan models.CarePlan
	if result := config.Db.Where("plant_id = ?", plant.ID).First(&carePlan); result.Error != nil {
		http.Error(w, "care plan not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(carePlan)
}

// UpdateCarePlan godoc
// @Summary      Update care plan
// @Description  Update the care plan for a plant
// @Tags         care-plan
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Plant ID"
// @Param        body body UpdateCarePlanRequest true "Update payload"
// @Success      200 {object} models.CarePlan
// @Failure      400 {string} string "invalid request"
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "care plan not found"
// @Router       /plants/{id}/care-plan [put]
func UpdateCarePlan(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		http.Error(w, "plant not found", http.StatusNotFound)
		return
	}

	var carePlan models.CarePlan
	if result := config.Db.Where("plant_id = ?", plant.ID).First(&carePlan); result.Error != nil {
		http.Error(w, "care plan not found", http.StatusNotFound)
		return
	}

	var req UpdateCarePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updates := map[string]interface{}{}
	if req.WateringFrequencyDays != nil {
		updates["watering_frequency_days"] = *req.WateringFrequencyDays
	}
	if req.WateringAmount != nil {
		updates["watering_amount"] = *req.WateringAmount
	}
	if req.WateringMethod != nil {
		updates["watering_method"] = *req.WateringMethod
	}
	if req.WateringTips != nil {
		updates["watering_tips"] = *req.WateringTips
	}
	if req.LightRequirement != nil {
		updates["light_requirement"] = *req.LightRequirement
	}
	if req.HumidityRequirement != nil {
		updates["humidity_requirement"] = *req.HumidityRequirement
	}

	if len(updates) > 0 {
		config.Db.Model(&carePlan).Updates(updates)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(carePlan)
}
