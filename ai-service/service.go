package aiservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Message represents a single chat message with its role and text content.
type Message struct {
	Role string `json:"role"` // "user" or "model"
	Text string `json:"text"`
}

// ChatRequest represents the JSON structure for an incoming chat message from the frontend.
// It now includes the previous chat history.
type ChatRequest struct {
	History        []Message `json:"history"`        // Array of previous messages
	CurrentMessage string    `json:"currentMessage"` // The new message from the user
}

// ChatResponse represents the JSON structure for an outgoing AI response to the frontend.
type ChatResponse struct {
	AIMessage         string `json:"aiMessage"`
	Error             string `json:"error,omitempty"`
	RawGeminiResponse string `json:"rawGeminiResponse,omitempty"`
}

// GeminiAPIRequest represents the payload structure for the Gemini API.
// Its Contents field will now hold the full conversation history.
type GeminiAPIRequest struct {
	Contents []struct {
		Role  string `json:"role"`
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
	GenerationConfig struct {
		ResponseMimeType string `json:"responseMimeType"`
	} `json:"generationConfig"`
}

// GeminiAPIResponse represents the response structure from the Gemini API.
type GeminiAPIResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	// Add other fields if you expect them, e.g.,
	// PromptFeedback struct {
	//     BlockReason string `json:"blockReason"`
	// } `json:"promptFeedback"`
}

func CallGemini(geminiReqPayload *GeminiAPIRequest) *ChatResponse {
	// var geminiReqPayload *GeminiAPIRequest
	// var messages []history.History
	// geminiReqPayload = InsertpreviousMessages(messages, currMessage)
	payloadBytes, err := json.Marshal(geminiReqPayload)
	if err != nil {
		return &ChatResponse{Error: fmt.Sprintf("Failed to marshal Gemini request: %v", err)}
	}

	// Log the outgoing request payload in a pretty JSON format
	var prettyReq bytes.Buffer
	json.Indent(&prettyReq, payloadBytes, "", "  ")
	// log.Printf("Gemini API Request Payload:\n%s", prettyReq.String())

	// Gemini API endpoint (gemini-2.0-flash)
	apiKey := "AIzaSyAsaFV-EnAtERnidP0rXxGs_mSI-JlBvjc" // Leave as empty string for Canvas to inject
	apiUrl := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	// Create HTTP client and request
	client := &http.Client{Timeout: 30 * time.Second}
	httpReq, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return &ChatResponse{Error: fmt.Sprintf("Failed to create HTTP request: %v", err)}
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Send request to Gemini API
	resp, err := client.Do(httpReq.WithContext(context.Background())) // Use context for request
	if err != nil {
		return &ChatResponse{Error: fmt.Sprintf("Failed to call Gemini API: %v", err)}
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ChatResponse{Error: fmt.Sprintf("Failed to read Gemini API response: %v", err)}
	}

	// Log the raw Gemini API response body in a pretty JSON format (for backend console)
	var prettyResp bytes.Buffer
	var rawGeminiResponseForFrontend string
	if err := json.Indent(&prettyResp, respBody, "", "  "); err != nil {
		log.Printf("Failed to pretty print Gemini API raw response: %v", err)
		// log.Printf("Gemini API Raw Response (unformatted):\n%s", string(respBody))
		rawGeminiResponseForFrontend = string(respBody)
	} else {
		// log.Printf("Gemini API Raw Response:\n%s", prettyResp.String())
		rawGeminiResponseForFrontend = prettyResp.String()
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Gemini API error response (%d): %s", resp.StatusCode, string(respBody))
		return &ChatResponse{Error: fmt.Sprintf("Gemini API returned error: %s", string(respBody)), RawGeminiResponse: rawGeminiResponseForFrontend}
	}

	// Parse Gemini API response
	var geminiResp GeminiAPIResponse
	if err := json.Unmarshal(respBody, &geminiResp); err != nil {
		return &ChatResponse{Error: fmt.Sprintf("Failed to parse Gemini API response: %v", err), RawGeminiResponse: rawGeminiResponseForFrontend}
	}

	// Log the unmarshaled Go struct for the Gemini response (for backend console)
	// log.Printf("Gemini API Response Struct: %+v", geminiResp)

	aiMessage := "Sorry, I couldn't generate a response."
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		aiMessage = geminiResp.Candidates[0].Content.Parts[0].Text
	}
	return &ChatResponse{
		AIMessage:         aiMessage,
		RawGeminiResponse: rawGeminiResponseForFrontend,
	}
}
