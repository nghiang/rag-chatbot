package models

import (
	"time"
	// "errors"
	// "gorm.io/gorm"
	"backend/internal/database"
)

type ChatMessage struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    SessionID uint      `gorm:"index;not null" json:"session_id"`
    Role      string    `gorm:"size:20;not null" json:"role"` // "user" | "assistant" | "system"
    Message   string    `gorm:"type:text" json:"message"`
    CreatedAt time.Time `json:"created_at"`
}

func MigrateChatMessage() error {
	return database.DB.AutoMigrate(&ChatMessage{})
}
