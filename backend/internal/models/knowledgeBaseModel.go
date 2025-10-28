package models

import (
	"time"
	// "errors"
	// "gorm.io/gorm"
	"backend/internal/database"
)

type KnowledgeBase struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    UserID      uint      `gorm:"index;not null" json:"user_id"`
    Name        string    `gorm:"size:255;not null" json:"name"`
    Description string    `json:"description,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    User       User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
    Documents  []Document `gorm:"foreignKey:KnowledgeBaseID" json:"documents,omitempty"`
    ChatSessions   []ChatSession  `gorm:"foreignKey:KnowledgeBaseID" json:"chat_sessions,omitempty"`
}

func MigrateKnowledgeBase() error {
	return database.DB.AutoMigrate(&KnowledgeBase{})
}
