package entiremedicalreport

import (
	"fmt"
	"log"
	"strconv" // Added for string manipulation

	// Import time for proper timestamps
	"github.com/gin-gonic/gin"
	aiservice "github.com/sukvij/health-checker/backend/ai-service"
	healthreportservice "github.com/sukvij/health-checker/backend/health-report-service"
	"gorm.io/gorm"
)

type GeminiAPIRequest struct {
	Contents         []ContentBlock   `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig"`
}

type ContentBlock struct {
	Role  string        `json:"role"` // "user", "system", or "assistant"
	Parts []ContentPart `json:"parts"`
}

type ContentPart struct {
	Text string `json:"text"`
}

type GenerationConfig struct {
	ResponseMimeType string `json:"responseMimeType"`
}

type HistoryService struct {
	Db *gorm.DB
}

func EntireMedicalReportRoute(app *gin.Engine, db *gorm.DB) {
	service := &HistoryService{Db: db}
	app.GET("/entire-report/:user_id", service.getEntireReportById)
}

func (service *HistoryService) getEntireReportById(ctx *gin.Context) {
	x := ctx.Param("user_id")
	userId, _ := strconv.Atoi(x)
	var previousReport []healthreportservice.HealthReport
	// Fetch previous messages for the user, excluding the current message being sent
	// Order by timestamp to maintain conversation flow
	err := service.Db.Where("user_id = ?", userId).Order("updated_at desc").Find(&previousReport).Error
	if err != nil {
		log.Printf("Error fetching previous reports for user %d: %v", userId, err)
		ctx.JSON(500, err)
		return
	}

	geminiReqPayload := InsertpreviousMessages(previousReport)

	result := aiservice.CallGemini(geminiReqPayload)
	ctx.JSON(200, result.AIMessage)
}

// InsertpreviousMessages prepares the payload for a regular chat conversation with Gemini.
func InsertpreviousMessages(previousReports []healthreportservice.HealthReport) *aiservice.GeminiAPIRequest {
	systemPrompt := `Below are multiple health logs from a user collected over time. Each log contains a date, type, and a description. 
Please analyze these logs and summarize the user's overall health condition till now. 
Mention common patterns, improvements, regressions, and any warning signs in simple, clear language. 
Conclude with any suggestions you would give based on these logs.`

	var geminiContents []struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}

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
			{Text: systemPrompt},
		},
	})

	// Append all user health logs
	for _, report := range previousReports {
		text := fmt.Sprintf("Date: %s\nType: %s\nDescription: %s", report.Date, report.Type, report.Description)
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
				{Text: text},
			},
		})
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
