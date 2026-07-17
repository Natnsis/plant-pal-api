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
	"plantPal/internals/response"
	"plantPal/internals/services"

	"github.com/gorilla/mux"
)

// StartDiagnosis godoc
// @Summary      Start a plant diagnosis
// @Description  Upload a plant image for AI-powered disease diagnosis and open a chat session
// @Tags         diagnosis
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        image  formData  file  true  "Plant image"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  response.ErrorResponse
// @Failure      500    {object}  response.ErrorResponse
// @Router       /diagnosis [post]
func StartDiagnosis(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "image file is required")
		return
	}
	defer file.Close()

	imageURL, err := services.UploadImage(file, fileHeader)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("failed to upload image: %s", err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	diagnosis, err := services.DiagnosePlant(ctx, imageURL)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("failed to diagnose plant: %s", err.Error()))
		return
	}

	session := models.AiChatSession{
		UserID: userID,
		Status: models.ChatStatusActive,
	}
	config.Db.Create(&session)

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

	scan := models.Scan{
		UserID:           userID,
		CapturedImageUrl: imageURL,
	}
	config.Db.Create(&scan)

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

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"session_id":   session.ID,
		"image_url":    imageURL,
		"diagnosis":    diagnosis,
		"chat_history": []models.AiChat{chat},
		"scan_id":      scan.ID,
	})
}

// GetDiagnosisChat godoc
// @Summary      Get diagnosis chat history
// @Description  Retrieve the chat history for a diagnosis session
// @Tags         diagnosis
// @Produce      json
// @Security     BearerAuth
// @Param        session_id  path  string  true  "Session ID"
// @Success      200         {object}  map[string]interface{}
// @Failure      404         {object}  response.ErrorResponse
// @Router       /diagnosis/{session_id} [get]
func GetDiagnosisChat(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	sessionID := mux.Vars(r)["session_id"]

	var session models.AiChatSession
	if result := config.Db.Where("id = ? AND user_id = ?", sessionID, userID).First(&session); result.Error != nil {
		response.Error(w, http.StatusNotFound, "session not found")
		return
	}

	var chats []models.AiChat
	config.Db.Where("session_id = ?", session.ID).Order("created_at asc").Find(&chats)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"session_id":   session.ID,
		"status":       session.Status,
		"chat_history": chats,
	})
}

// SendChatMessage godoc
// @Summary      Send a chat message
// @Description  Send a follow-up message in a diagnosis chat session
// @Tags         diagnosis
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        session_id  path      string              true  "Session ID"
// @Param        body        body      ChatMessageRequest  true  "Chat message"
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  response.ErrorResponse
// @Failure      404         {object}  response.ErrorResponse
// @Router       /diagnosis/{session_id}/chat [post]
func SendChatMessage(w http.ResponseWriter, r *http.Request) {
	userID := middlewares.GetUserID(r)
	sessionID := mux.Vars(r)["session_id"]

	var session models.AiChatSession
	if result := config.Db.Where("id = ? AND user_id = ? AND status = ?", sessionID, userID, models.ChatStatusActive).First(&session); result.Error != nil {
		response.Error(w, http.StatusNotFound, "session not found or archived")
		return
	}

	var req ChatMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Message == "" {
		response.Error(w, http.StatusBadRequest, "message is required")
		return
	}

	userChat := models.AiChat{
		SessionID:   session.ID,
		SenderType:  models.SenderTypeUser,
		MessageBody: req.Message,
	}
	config.Db.Create(&userChat)

	var allChats []models.AiChat
	config.Db.Where("session_id = ?", session.ID).Order("created_at asc").Find(&allChats)

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
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("failed to get AI response: %s", err.Error()))
		return
	}

	aiChat := models.AiChat{
		SessionID:   session.ID,
		SenderType:  models.SenderTypeAI,
		MessageBody: aiResponse,
	}
	config.Db.Create(&aiChat)

	var updatedChats []models.AiChat
	config.Db.Where("session_id = ?", session.ID).Order("created_at asc").Find(&updatedChats)

	response.JSON(w, http.StatusOK, map[string]interface{}{
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
