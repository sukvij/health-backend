package chatservice

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings" // Added for string manipulation

	"github.com/gin-gonic/gin"
	aiservice "github.com/sukvij/health-checker/backend/ai-service"
	"gorm.io/gorm"
)

type HistoryService struct {
	Db *gorm.DB
}

func HistroyRoute(app *gin.Engine, db *gorm.DB) {
	service := &HistoryService{Db: db}
	app.GET("/history/:user_id", service.GetAllHistoryById)
	app.POST("/history", service.CreateHistory)
}

func (service *HistoryService) CreateHistory(ctx *gin.Context) {
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
	isReportRequest := strings.Contains(strings.ToLower(chat.Text), "generate medical report") ||
		strings.Contains(strings.ToLower(chat.Text), "medical summary")

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
	// It's good practice to set the current timestamp before saving.
	// You might want to use time.Now().Format(...) for a proper timestamp.
	// For simplicity, I'm just showing a placeholder if not already handled by GORM's created_at.
	// chat.Timestamp = time.Now().Format(time.RFC3339) // Example

	if err := service.Db.Create(&chat).Error; err != nil {
		log.Printf("Error saving chat message to DB: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chat message."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"response": chat.Response, "user_message": chat.Text})
}

func (service *HistoryService) GetAllHistoryById(ctx *gin.Context) {
	var conversationSummaries []ChatMessage
	x := ctx.Param("user_id")
	userId, err := strconv.Atoi(x)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = service.Db.Where("user_id = ?", userId).Order("timestamp asc").Find(&conversationSummaries).Error // Order by timestamp
	if err != nil {
		log.Printf("Error fetching all conversation summaries: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversation list"})
		return
	}

	ctx.JSON(http.StatusOK, conversationSummaries)
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
	for _, msg := range previousMessages {
		fullConversation += fmt.Sprintf("User: %s\n", msg.Text)
		if msg.Response != "" {
			fullConversation += fmt.Sprintf("AI: %s\n", msg.Response)
		}
	}

	// This is the crucial part: craft a clear instruction for Gemini.
	prompt := fmt.Sprintf(`Based on the following conversation history, please generate a concise medical report or summary. Focus on key symptoms, conditions discussed, medical advice given, and any other relevant health information. Format it as a clear, short report (under 200 words if possible), without conversational elements.

Conversation History:
%s

Medical Report:`, fullConversation)

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
