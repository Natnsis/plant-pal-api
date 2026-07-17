package handlers

import (
	"encoding/json"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
	"plantPal/internals/models"
	"plantPal/internals/response"

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

func GetCarePlan(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		response.Error(w, http.StatusNotFound, "plant not found")
		return
	}

	var carePlan models.CarePlan
	if result := config.Db.Where("plant_id = ?", plant.ID).First(&carePlan); result.Error != nil {
		response.Error(w, http.StatusNotFound, "care plan not found")
		return
	}

	response.JSON(w, http.StatusOK, carePlan)
}

func UpdateCarePlan(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		response.Error(w, http.StatusNotFound, "plant not found")
		return
	}

	var carePlan models.CarePlan
	if result := config.Db.Where("plant_id = ?", plant.ID).First(&carePlan); result.Error != nil {
		response.Error(w, http.StatusNotFound, "care plan not found")
		return
	}

	var req UpdateCarePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
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

	response.JSON(w, http.StatusOK, carePlan)
}
