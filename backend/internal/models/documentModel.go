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
    Description     string    `gorm:"size:255" json:"description"`
    EmbeddingStatus string    `gorm:"size:20;default:'pending'" json:"embedding_status"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`

    KnowledgeBase KnowledgeBase `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"knowledge_base,omitempty"`
    User          User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

func MigrateDocument() error {
	return database.DB.AutoMigrate(&Document{})
}

func CreateDocument(d *Document) error {
    return database.DB.Create(d).Error
}

func GetDocumentByID(id uint) (*Document, error) {
    var doc Document
    if err := database.DB.First(&doc, id).Error; err != nil {
        return nil, err
    }
    return &doc, nil
}

func UpdateDocument(d *Document) error {
    return database.DB.Save(d).Error
}

func UpdateDocumentEmbeddingStatus(id uint, status string) error {
    return database.DB.Model(&Document{}).Where("id = ?", id).Update("embedding_status", status).Error
}
