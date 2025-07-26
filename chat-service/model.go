package chatservice

// ChatMessage represents a single message within a chat session.
// This is a helper struct for JSON serialization within ChatSession.
// It does NOT correspond to a direct database table.
// Corresponds to frontend/app/types/chat.ts (ChatMessage interface)
type ChatMessage struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    uint   `gorm:"user_id" json:"user_id"`
	Text      string `json:"text"`
	Sender    string `json:"sender"`    // "user" or "ai"
	Timestamp string `json:"timestamp"` // ISO string
	Response  string `json:"response"`
}
