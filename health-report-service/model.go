package healthreportservice

import (
	"time"

	"gorm.io/gorm"
)

// HealthReport represents a health log submitted by a user.
// Corresponds to frontend/app/types/health_report.ts
type HealthReport struct {
	// ID is the primary key for the HealthReport, now auto-incrementing integer.
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	UserID uint `gorm:"not null;index" json:"userId"` // References User's uint ID

	// Type of the report (e.g., "Daily Log", "Symptom Update").
	Type string `gorm:"type:varchar(100);not null" json:"type"`

	// Description of the health report. Uses TEXT type for potentially long content.
	Description string `gorm:"type:text;not null" json:"description"`

	// Date the report refers to (e.g., "2023-11-01"). Stored as string for simplicity as per TS.
	Date string `gorm:"type:varchar(10);not null" json:"date"` // Format: YYYY-MM-DD

	// GORM's default fields for timestamps and soft delete.
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:true" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
