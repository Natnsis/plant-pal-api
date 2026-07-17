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

func GetPlantReminders(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	plantID := mux.Vars(r)["id"]

	var plant models.Plant
	if result := config.Db.Where("id = ? AND user_id = ?", plantID, userID).First(&plant); result.Error != nil {
		response.Error(w, http.StatusNotFound, "plant not found")
		return
	}

	var reminders []models.Reminder
	config.Db.Where("plant_id = ?", plant.ID).Order("scheduled_time asc").Find(&reminders)

	response.JSON(w, http.StatusOK, reminders)
}

func GetTodayReminders(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	var reminders []models.Reminder
	config.Db.Joins("JOIN plants ON plants.id = reminders.plant_id").
		Where("plants.user_id = ? AND reminders.scheduled_time >= ? AND reminders.scheduled_time < ? AND reminders.is_completed = false",
			userID, startOfDay, endOfDay).
		Order("reminders.scheduled_time asc").
		Find(&reminders)

	response.JSON(w, http.StatusOK, reminders)
}

type UpdateReminderRequest struct {
	IsCompleted *bool `json:"is_completed"`
	Snooze      bool  `json:"snooze"`
}

func UpdateReminder(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	reminderID := mux.Vars(r)["id"]

	var reminder models.Reminder
	if result := config.Db.Joins("JOIN plants ON plants.id = reminders.plant_id").
		Where("reminders.id = ? AND plants.user_id = ?", reminderID, userID).
		First(&reminder); result.Error != nil {
		response.Error(w, http.StatusNotFound, "reminder not found")
		return
	}

	var req UpdateReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.IsCompleted != nil && *req.IsCompleted {
		reminder.IsCompleted = true
		reminder.CompletedAt = time.Now()
		config.Db.Save(&reminder)

		var plant models.Plant
		config.Db.Where("id = ?", reminder.PlantID).First(&plant)
		config.Db.Model(&models.User{}).Where("id = ?", plant.UserID).
			UpdateColumn("total_task_done", 1)
	} else if req.Snooze {
		reminder.SnoozeCount++
		reminder.ScheduledTime = reminder.ScheduledTime.Add(15 * time.Minute)
		config.Db.Save(&reminder)
	}

	response.JSON(w, http.StatusOK, reminder)
}
