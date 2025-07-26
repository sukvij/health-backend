package userservice

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	// ID is the primary key for the User, now auto-incrementing integer.
	// `gorm:"primaryKey;autoIncrement"`: Specifies ID as primary key and auto-incrementing.
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"` // Use uint for auto-incrementing IDs

	// Name of the user.
	Name string `gorm:"type:varchar(255);not null" json:"name"`

	// Email of the user, must be unique.
	Email string `gorm:"type:varchar(255);unique;not null" json:"email"`

	// GORM's default fields for creation, update timestamps, and soft delete.
	// GORM automatically populates these if their names are exactly CreatedAt, UpdatedAt, DeletedAt.
	// `autoUpdateTime:true` ensures UpdatedAt is updated on every save.
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:true" json:"updatedAt"` // MySQL typically uses DATETIME, autoUpdateTime handles this
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                       // Soft delete, `json:"-"` hides it from JSON output
}
