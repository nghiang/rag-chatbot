package models

import (
	"time"
	// "errors"
	// "gorm.io/gorm"
	"backend/internal/database"
)


type Document struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    KnowledgeBaseID uint      `gorm:"index" json:"knowledge_base_id"`
    UserID          uint      `gorm:"index;not null" json:"user_id"`
    Name            string    `gorm:"size:255" json:"name"`
    FileType        string    `gorm:"size:50" json:"file_type"`
    MinioObjectName string    `gorm:"size:255" json:"minio_object_name"`
    EmbeddingStatus string    `gorm:"size:20;default:'pending'" json:"embedding_status"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`

    KnowledgeBase KnowledgeBase `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"knowledge_base,omitempty"`
    User          User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

func MigrateDocument() error {
	return database.DB.AutoMigrate(&Document{})
}
