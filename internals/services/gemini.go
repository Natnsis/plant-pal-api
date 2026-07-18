package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"plantPal/internals/config"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PlantIdentification struct {
	CommonName       string   `json:"common_name"`
	ScientificName   string   `json:"scientific_name"`
	Family           string   `json:"family"`
	Origin           string   `json:"origin"`
	ConfidenceScore  float64  `json:"confidence_score"`
	HealthAssessment string   `json:"health_assessment"`
	DetectedSymptoms []string `json:"detected_symptoms"`
	PrimaryAssessment string  `json:"primary_assessment"`
	TreatmentSteps   []string `json:"treatment_steps"`
	CareRecommendations CareRecommendations `json:"care_recommendations"`
}

type CareRecommendations struct {
	WateringFrequencyDays uint   `json:"watering_frequency_days"`
	WateringAmount        string `json:"watering_amount"`
	WateringMethod        string `json:"watering_method"`
	WateringTips          string `json:"watering_tips"`
	LightRequirement      string `json:"light_requirement"`
	HumidityRequirement   string `json:"humidity_requirement"`
}

type DiagnosisResult struct {
	PlantType      string   `json:"plant_type"`
	IssueDescription string `json:"issue_description"`
	Severity       string   `json:"severity"`
	Causes         []string `json:"causes"`
	Solutions      []string `json:"solutions"`
	PreventionTips []string `json:"prevention_tips"`
}

func newGeminiClient(ctx context.Context) (*genai.Client, error) {
	return genai.NewClient(ctx, option.WithAPIKey(config.GeminiAPIKey))
}

func IdentifyPlant(ctx context.Context, imageURL string) (*PlantIdentification, error) {
	client, err := newGeminiClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	defer client.Close()

	imgData, err := downloadImage(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash")
	model.ResponseMIMEType = "application/json"

	prompt := `You are a plant identification expert. Analyze this plant image and return a JSON object with the following fields:
{
  "common_name": "string - common name of the plant",
  "scientific_name": "string - scientific/binomial name",
  "family": "string - plant family",
  "origin": "string - geographic origin",
  "confidence_score": number between 0 and 1 (how confident you are in the identification),
  "health_assessment": "string - overall health description",
  "detected_symptoms": ["string - any visible symptoms like yellow leaves, brown spots, etc"],
  "primary_assessment": "string - main health issue or status",
  "treatment_steps": ["string - step by step treatment if issues found"],
  "care_recommendations": {
    "watering_frequency_days": number (e.g. 7),
    "watering_amount": "string (e.g. 200ml)",
    "watering_method": "string (e.g. top watering, bottom watering)",
    "watering_tips": "string",
    "light_requirement": "string (e.g. bright indirect light)",
    "humidity_requirement": "string (e.g. 60-80%)"
  }
}

Return ONLY valid JSON, no markdown code fences.`

	resp, err := model.GenerateContent(ctx,
		genai.ImageData("image/jpeg", imgData),
		genai.Text(prompt),
	)
	if err != nil {
		return nil, fmt.Errorf("gemini generate failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from gemini")
	}

	textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response type from gemini")
	}
	text := cleanJSONResponse(string(textPart))

	var result PlantIdentification
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %w", err)
	}

	return &result, nil
}

func DiagnosePlant(ctx context.Context, imageURL string) (*DiagnosisResult, error) {
	client, err := newGeminiClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	defer client.Close()

	imgData, err := downloadImage(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash")
	model.ResponseMIMEType = "application/json"

	prompt := `You are a plant diagnosis expert. Analyze this plant image and diagnose any issues.
Return a JSON object:
{
  "plant_type": "string - type/name of the plant",
  "issue_description": "string - detailed description of the problem",
  "severity": "string - low/medium/high/critical",
  "causes": ["string - possible causes"],
  "solutions": ["string - step by step solutions"],
  "prevention_tips": ["string - how to prevent this in future"]
}

Return ONLY valid JSON, no markdown code fences.`

	resp, err := model.GenerateContent(ctx,
		genai.ImageData("image/jpeg", imgData),
		genai.Text(prompt),
	)
	if err != nil {
		return nil, fmt.Errorf("gemini generate failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from gemini")
	}

	textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return nil, fmt.Errorf("unexpected response type from gemini")
	}
	text := cleanJSONResponse(string(textPart))

	var result DiagnosisResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %w", err)
	}

	return &result, nil
}

func ChatWithAI(ctx context.Context, history []ChatMessage, newMessage string) (string, error) {
	client, err := newGeminiClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create gemini client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")

	var chatHistory []*genai.Content
	for _, msg := range history {
		role := "user"
		if msg.Role == "ai" {
			role = "model"
		}
		chatHistory = append(chatHistory, &genai.Content{
			Role:  role,
			Parts: []genai.Part{genai.Text(msg.Content)},
		})
	}

	session := model.StartChat()
	session.History = chatHistory

	resp, err := session.SendMessage(ctx, genai.Text(newMessage))
	if err != nil {
		return "", fmt.Errorf("gemini chat failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from gemini")
	}

	textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", fmt.Errorf("unexpected response type from gemini")
	}

	return string(textPart), nil
}

func downloadImage(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func cleanJSONResponse(text string) string {
	if len(text) >= 3 && text[:3] == "```" {
		start := 0
		for i := 0; i < len(text); i++ {
			if text[i] == '\n' {
				start = i + 1
				break
			}
		}
		end := len(text)
		if len(text) >= 3 && text[end-3:] == "```" {
			end -= 3
		}
		text = text[start:end]
	}
	return text
}
