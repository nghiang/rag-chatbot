package models

import (
	"time"
	// "errors"
	// "gorm.io/gorm"
	// "backend/internal/database"
)

type ChatSession struct {
    ID              uint       `gorm:"primaryKey" json:"id"`
    UserID          uint       `gorm:"index;not null" json:"user_id"`
    KnowledgeBaseID *uint      `gorm:"index" json:"knowledge_base_id,omitempty"`
    Title           string     `gorm:"size:255" json:"title"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`

    User          User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
    KnowledgeBase *KnowledgeBase `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"knowledge_base,omitempty"`
    Messages      []ChatMessage `gorm:"foreignKey:SessionID" json:"messages,omitempty"`
}


