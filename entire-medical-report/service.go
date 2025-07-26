package entiremedicalreport

import (
	"fmt"
	"log"
	"net/http"
	"strings" // Added for string manipulation
	"time"    // Import time for proper timestamps

	"github.com/gin-gonic/gin"
	aiservice "github.com/sukvij/health-checker/backend/ai-service"
	"gorm.io/gorm"
)

// ChatMessage struct provided by the user.
// Using it as is, but typically, 'Timestamp' would be a time.Time type in Go.
type ChatMessage struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    uint   `gorm:"user_id" json:"user_id"`
	Text      string `json:"text"`
	Sender    string `json:"sender"`    // "user" or "ai"
	Timestamp string `json:"timestamp"` // ISO string
	Response  string `json:"response"`
}

type HistoryService struct {
	Db *gorm.DB
}

func EntireMedicalReportRoute(app *gin.Engine, db *gorm.DB) {
	service := &HistoryService{Db: db}
	app.GET("/entire-report/:user_id", service.GetAllHistoryById)
}

func (service *HistoryService) GetAllHistoryById(ctx *gin.Context) {
	var chat ChatMessage
	if err := ctx.ShouldBindJSON(&chat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request payload: %v", err)})
		return
	}

	var previousMessages []ChatMessage
	// Fetch previous messages for the user, excluding the current message being sent
	// Order by timestamp to maintain conversation flow
	err := service.Db.Where("user_id = ?", chat.UserId).Order("timestamp asc").Find(&previousMessages).Error
	if err != nil {
		log.Printf("Error fetching previous messages for user %d: %v", chat.UserId, err)
		// Even if there's an error, we can proceed with an empty history
	}

	// Determine if the current message is a trigger for a medical report
	// Make the trigger phrases more specific to avoid accidental triggers
	isReportRequest := strings.Contains(strings.ToLower(chat.Text), "medical report") ||
		strings.Contains(strings.ToLower(chat.Text), "health summary") ||
		strings.Contains(strings.ToLower(chat.Text), "generate report") ||
		strings.Contains(strings.ToLower(chat.Text), "meri aaj tak ki previous history dekh ke short report bnao") // Added your specific trigger

	var geminiReqPayload *aiservice.GeminiAPIRequest

	if isReportRequest {
		log.Printf("Medical report requested by user %d.", chat.UserId)
		// If a report is requested, generate a prompt specifically for the medical report
		geminiReqPayload = CreateMedicalReportPayload(previousMessages)
	} else {
		// Otherwise, continue with regular chat conversation
		geminiReqPayload = InsertpreviousMessages(previousMessages, chat.Text)
	}

	result := aiservice.CallGemini(geminiReqPayload)

	if result.Error != "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	chat.Response = result.AIMessage
	// Set the current timestamp before saving.
	chat.Timestamp = time.Now().Format(time.RFC3339) // Use RFC3339 for ISO string

	if err := service.Db.Create(&chat).Error; err != nil {
		log.Printf("Error saving chat message to DB: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chat message."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": chat.Response, "user_message": chat.Text})
}

// InsertpreviousMessages prepares the payload for a regular chat conversation with Gemini.
func InsertpreviousMessages(previousMessages []ChatMessage, currMessage string) *aiservice.GeminiAPIRequest {
	var geminiContents []struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}

	// Add previous messages from history to Gemini's contents
	for _, msg := range previousMessages {
		// User's message
		geminiContents = append(geminiContents, struct {
			Role  string `json:"role"`
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			Role: "user",
			Parts: []struct {
				Text string `json:"text"`
			}{
				{Text: msg.Text},
			},
		})
		// AI's response to the user's message
		if msg.Response != "" {
			geminiContents = append(geminiContents, struct {
				Role  string `json:"role"`
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			}{
				Role: "model", // Assuming "model" is the AI's role for responses
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: msg.Response},
				},
			})
		}
	}

	// Add the current user message
	geminiContents = append(geminiContents, struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{
		Role: "user",
		Parts: []struct {
			Text string `json:"text"`
		}{
			{Text: currMessage},
		},
	})

	return &aiservice.GeminiAPIRequest{
		Contents: geminiContents,
		GenerationConfig: struct {
			ResponseMimeType string `json:"responseMimeType"`
		}{
			ResponseMimeType: "text/plain",
		},
	}
}

// CreateMedicalReportPayload creates a specific payload for generating a medical report.
func CreateMedicalReportPayload(previousMessages []ChatMessage) *aiservice.GeminiAPIRequest {
	var fullConversation string
	if len(previousMessages) == 0 {
		fullConversation = "No prior conversation history available."
	} else {
		for _, msg := range previousMessages {
			fullConversation += fmt.Sprintf("User: %s\n", msg.Text)
			if msg.Response != "" {
				fullConversation += fmt.Sprintf("AI: %s\n", msg.Response)
			}
		}
	}

	// This is the crucial part: craft a clear instruction for Gemini.
	// Emphasize that Gemini should only use the provided history.
	prompt := fmt.Sprintf(`You are an AI medical assistant tasked with summarizing a user's health history.
Based *solely* on the following conversation history, please generate a concise medical report or summary.
Do NOT access any external information or personal data not explicitly provided in the conversation.
Focus on extracting and summarizing key symptoms, diagnosed conditions, medical advice given by the AI, and any other relevant health information mentioned by the user or AI.
Present the information as a clear, short report (aim for under 200 words if possible), using bullet points or a structured paragraph format, without any conversational preamble or disclaimers about being an AI or not having memory. Just provide the report.

Conversation History:
---
%s
---

Medical Report Summary:`, fullConversation)

	geminiContents := []struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{
		{
			Role: "user",
			Parts: []struct {
				Text string `json:"text"`
			}{
				{Text: prompt},
			},
		},
	}

	return &aiservice.GeminiAPIRequest{
		Contents: geminiContents,
		GenerationConfig: struct {
			ResponseMimeType string `json:"responseMimeType"`
		}{
			ResponseMimeType: "text/plain",
		},
	}
}
