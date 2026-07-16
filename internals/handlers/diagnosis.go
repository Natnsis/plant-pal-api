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

// StartDiagnosis godoc
// @Summary      Start a diagnosis chat
// @Description  Upload a plant image to start a diagnosis chat session. Returns initial analysis and chat session.
// @Tags         diagnosis
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        image formData file true "Plant image"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {string} string "missing image"
// @Failure      401 {string} string "unauthorized"
// @Failure      500 {string} string "internal error"
// @Router       /diagnosis [post]
func StartDiagnosis(w http.ResponseWriter, r *http.Request) {
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

	diagnosis, err := services.DiagnosePlant(ctx, imageURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to diagnose plant: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Create chat session
	session := models.AiChatSession{
		UserID: userID,
		Status: models.ChatStatusActive,
	}
	config.Db.Create(&session)

	// Format initial AI message with diagnosis
	aiMessage := fmt.Sprintf(
		"Plant Type: %s\n\nIssue: %s\n\nSeverity: %s\n\nCauses:\n%s\n\nSolutions:\n%s\n\nPrevention Tips:\n%s",
		diagnosis.PlantType,
		diagnosis.IssueDescription,
		diagnosis.Severity,
		joinBullets(diagnosis.Causes),
		joinBullets(diagnosis.Solutions),
		joinBullets(diagnosis.PreventionTips),
	)

	chat := models.AiChat{
		SessionID:   session.ID,
		SenderType:  models.SenderTypeAI,
		MessageBody: aiMessage,
	}
	config.Db.Create(&chat)

	// Create scan record for the image
	scan := models.Scan{
		UserID:           userID,
		CapturedImageUrl: imageURL,
	}
	config.Db.Create(&scan)

	// Create analysis result
	analysis := models.AiAnalysisResult{
		ScanID:             scan.ID,
		AiModelVersion:     "gemini-1.5-flash",
		ConfidenceScore:    1.0,
		AnalysisType:       models.DiagnosisAnalysisType,
		PrimaryAssessment:  diagnosis.IssueDescription,
		TreatmentPlanSteps: joinStrings(diagnosis.Solutions),
		MetadataPayload:    toJSON(diagnosis),
	}
	config.Db.Create(&analysis)

	scan.AnalysisID = analysis.ID
	config.Db.Save(&scan)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_id":        session.ID,
		"image_url":         imageURL,
		"diagnosis":         diagnosis,
		"chat_history":      []models.AiChat{chat},
		"scan_id":           scan.ID,
	})
}

// GetDiagnosisChat godoc
// @Summary      Get diagnosis chat history
// @Description  Get full chat history for a diagnosis session
// @Tags         diagnosis
// @Produce      json
// @Security     BearerAuth
// @Param        session_id path int true "Chat Session ID"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "session not found"
// @Router       /diagnosis/{session_id} [get]
func GetDiagnosisChat(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	sessionID := mux.Vars(r)["session_id"]

	var session models.AiChatSession
	if result := config.Db.Where("id = ? AND user_id = ?", sessionID, userID).First(&session); result.Error != nil {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	var chats []models.AiChat
	config.Db.Where("session_id = ?", session.ID).Order("created_at asc").Find(&chats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_id":  session.ID,
		"status":      session.Status,
		"chat_history": chats,
	})
}

// SendChatMessage godoc
// @Summary      Send a follow-up message in a diagnosis chat
// @Description  Send a message in an existing diagnosis session. Full chat history is sent to AI.
// @Tags         diagnosis
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        session_id path int true "Chat Session ID"
// @Param        body body ChatMessageRequest true "Message"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {string} string "invalid request"
// @Failure      401 {string} string "unauthorized"
// @Failure      404 {string} string "session not found"
// @Router       /diagnosis/{session_id}/chat [post]
func SendChatMessage(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	sessionID := mux.Vars(r)["session_id"]

	var session models.AiChatSession
	if result := config.Db.Where("id = ? AND user_id = ? AND status = ?", sessionID, userID, models.ChatStatusActive).First(&session); result.Error != nil {
		http.Error(w, "session not found or archived", http.StatusNotFound)
		return
	}

	var req ChatMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	// Save user message
	userChat := models.AiChat{
		SessionID:   session.ID,
		SenderType:  models.SenderTypeUser,
		MessageBody: req.Message,
	}
	config.Db.Create(&userChat)

	// Get full chat history
	var allChats []models.AiChat
	config.Db.Where("session_id = ?", session.ID).Order("created_at asc").Find(&allChats)

	// Convert to services format
	var history []services.ChatMessage
	for _, c := range allChats {
		history = append(history, services.ChatMessage{
			Role:    string(c.SenderType),
			Content: c.MessageBody,
		})
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	aiResponse, err := services.ChatWithAI(ctx, history, req.Message)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get AI response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Save AI response
	aiChat := models.AiChat{
		SessionID:   session.ID,
		SenderType:  models.SenderTypeAI,
		MessageBody: aiResponse,
	}
	config.Db.Create(&aiChat)

	// Reload all chats including new ones
	var updatedChats []models.AiChat
	config.Db.Where("session_id = ?", session.ID).Order("created_at asc").Find(&updatedChats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_id":   session.ID,
		"user_message": userChat,
		"ai_message":   aiChat,
		"chat_history": updatedChats,
	})
}

type ChatMessageRequest struct {
	Message string `json:"message"`
}

func joinBullets(items []string) string {
	result := ""
	for _, item := range items {
		result += "- " + item + "\n"
	}
	return result
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
