package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
	"plantPal/internals/models"
	"plantPal/internals/services"

	"github.com/gorilla/mux"
)

// ScanPlant godoc
// @Summary      Scan a plant
// @Description  Upload a plant image for identification. Returns retake=true if confidence is low.
// @Tags         scan
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        image formData file true "Plant image"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {string} string "missing image"
// @Failure      401 {string} string "unauthorized"
// @Failure      500 {string} string "internal error"
// @Router       /scan [post]
func ScanPlant(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "image file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageURL, err := services.UploadImage(file, fileHeader)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to upload image: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	identification, err := services.IdentifyPlant(ctx, imageURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to analyze plant: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Create scan record
	scan := models.Scan{
		UserID:           userID,
		CapturedImageUrl: imageURL,
		ConfidenceScore:  identification.ConfidenceScore,
	}

	// If confidence is low, return retake request
	if identification.ConfidenceScore < 0.7 {
		scan.Retake = true
		config.Db.Create(&scan)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"retake":            true,
			"confidence_score":  identification.ConfidenceScore,
			"message":           "Could not confidently identify this plant. Please retake with a clearer image.",
			"partial_data":      identification,
			"scan_id":           scan.ID,
			"captured_image_url": imageURL,
		})
		return
	}

	// Good confidence - create full analysis
	// Upsert species by scientific name
	var species models.Species
	result := config.Db.Where("scientific_name = ?", identification.ScientificName).First(&species)
	if result.Error != nil {
		species = models.Species{
			CommonName:      identification.CommonName,
			ScientificName:  identification.ScientificName,
			Family:          identification.Family,
			Origin:          identification.Origin,
			DifficultyLevel: models.MediumDifficulty,
		}
		config.Db.Create(&species)
	}

	// Create plant
	plant := models.Plant{
		UserID:      userID,
		SpeciesID:   &species.ID,
		Nickname:    identification.CommonName,
		HealthScore: int(identification.ConfidenceScore * 100),
		Status:      models.GoodPlantStatus,
	}
	config.Db.Create(&plant)

	// Update scan with plant ID
	scan.PlantID = plant.ID
	config.Db.Save(&scan)

	// Create analysis result
	analysis := models.AiAnalysisResult{
		ScanID:             scan.ID,
		AiModelVersion:     "gemini-1.5-flash",
		ConfidenceScore:    identification.ConfidenceScore,
		AnalysisType:       models.IdentificationAnalysisType,
		DetectedSymptoms:   models.StringList(identification.DetectedSymptoms),
		PrimaryAssessment:  identification.PrimaryAssessment,
		TreatmentPlanSteps: joinStrings(identification.TreatmentSteps),
	}
	config.Db.Create(&analysis)

	scan.AnalysisID = analysis.ID
	config.Db.Save(&scan)

	// Create care plan
	carePlan := models.CarePlan{
		PlantID:               plant.ID,
		WateringFrequencyDays: identification.CareRecommendations.WateringFrequencyDays,
		WateringAmount:        identification.CareRecommendations.WateringAmount,
		WateringMethod:        identification.CareRecommendations.WateringMethod,
		WateringTips:          identification.CareRecommendations.WateringTips,
		LightRequirement:      identification.CareRecommendations.LightRequirement,
		HumidityRequirement:   identification.CareRecommendations.HumidityRequirement,
	}
	config.Db.Create(&carePlan)

	// Generate reminders based on care plan
	now := time.Now()
	createReminderIfNotExists(plant.ID, models.WaterTask, now.AddDate(0, 0, int(identification.CareRecommendations.WateringFrequencyDays)))
	createReminderIfNotExists(plant.ID, models.FertilizeTask, now.AddDate(0, 1, 0))
	createReminderIfNotExists(plant.ID, models.RotateTask, now.AddDate(0, 0, 7))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"retake":           false,
		"scan_id":          scan.ID,
		"plant":            plant,
		"species":          species,
		"analysis":         analysis,
		"care_plan":        carePlan,
		"captured_image_url": imageURL,
	})
}

// GetScan godoc
// @Summary      Get scan details
// @Description  Get details of a specific scan
// @Tags         scan
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Scan ID"
// @Success      200 {object} models.Scan
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "scan not found"
// @Router       /scan/{id} [get]
func GetScan(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	scanID := mux.Vars(r)["id"]

	var scan models.Scan
	if result := config.Db.Where("id = ? AND user_id = ?", scanID, userID).
		Preload("Plant").Preload("Plant.Species").First(&scan); result.Error != nil {
		http.Error(w, "scan not found", http.StatusNotFound)
		return
	}

	var analysis models.AiAnalysisResult
	config.Db.Where("scan_id = ?", scan.ID).First(&analysis)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"scan":     scan,
		"analysis": analysis,
	})
}

func createReminderIfNotExists(plantID uint, taskType models.TaskType, scheduledTime time.Time) {
	var count int64
	config.Db.Model(&models.Reminder{}).
		Where("plant_id = ? AND task_type = ? AND is_completed = false", plantID, taskType).
		Count(&count)
	if count == 0 {
		reminder := models.Reminder{
			PlantID:       plantID,
			TaskType:      taskType,
			ScheduledTime: scheduledTime,
		}
		config.Db.Create(&reminder)
	}
}

func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += "\n"
		}
		result += s
	}
	return result
}
