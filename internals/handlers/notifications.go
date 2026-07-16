package handlers

import (
	"encoding/json"
	"net/http"

	"plantPal/internals/config"
	"plantPal/internals/middlewares"
	"plantPal/internals/models"
)

type UpdateNotificationRequest struct {
	NotificationEnabled         *bool   `json:"notification_enabled"`
	DailySummaryEnabled         *bool   `json:"daily_summary_enabled"`
	SoundAlertEnabled           *bool   `json:"sound_alert_enabled"`
	VibrationEnabled            *bool   `json:"vibration_enabled"`
	PreferredNotificationTime   *string `json:"preferred_notification_time"`
	DefaultSnoozeDurationMinute *uint   `json:"default_snooze_duration_minute"`
}

// GetNotificationSettings godoc
// @Summary      Get notification settings
// @Description  Get the notification preferences for the authenticated user
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.Notification
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "settings not found"
// @Router       /notifications [get]
func GetNotificationSettings(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	var notification models.Notification
	if result := config.Db.Where("user_id = ?", userID).First(&notification); result.Error != nil {
		// Create default settings if none exist
		notification = models.Notification{
			UserID:                userID,
			NotificationEnabled:   true,
			DailySummaryEnabled:   false,
			SoundAlertEnabled:     true,
			VibrationEnabled:      true,
		}
		config.Db.Create(&notification)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}

// UpdateNotificationSettings godoc
// @Summary      Update notification settings
// @Description  Update the notification preferences for the authenticated user
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body UpdateNotificationRequest true "Settings payload"
// @Success      200 {object} models.Notification
// @Failure      400 {string} string "invalid request"
// @Failure      401 {string} string "unauthorized"
// @Router       /notifications [put]
func UpdateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	var notification models.Notification
	if result := config.Db.Where("user_id = ?", userID).First(&notification); result.Error != nil {
		http.Error(w, "settings not found", http.StatusNotFound)
		return
	}

	var req UpdateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updates := map[string]interface{}{}
	if req.NotificationEnabled != nil {
		updates["notification_enabled"] = *req.NotificationEnabled
	}
	if req.DailySummaryEnabled != nil {
		updates["daily_summary_enabled"] = *req.DailySummaryEnabled
	}
	if req.SoundAlertEnabled != nil {
		updates["sound_alert_enabled"] = *req.SoundAlertEnabled
	}
	if req.VibrationEnabled != nil {
		updates["vibration_enabled"] = *req.VibrationEnabled
	}
	if req.DefaultSnoozeDurationMinute != nil {
		updates["default_snooze_duration_minute"] = *req.DefaultSnoozeDurationMinute
	}
	if req.PreferredNotificationTime != nil {
		updates["preferred_notification_time"] = *req.PreferredNotificationTime
	}

	if len(updates) > 0 {
		config.Db.Model(&notification).Updates(updates)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}
