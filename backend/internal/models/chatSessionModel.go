package models

import (
	"backend/internal/services"
	"time"
)

type ChatSession struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"index;not null" json:"user_id"`
	KnowledgeBaseID *uint     `gorm:"index" json:"knowledge_base_id,omitempty"`
	Title           string    `gorm:"size:255" json:"title"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	User          User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	KnowledgeBase *KnowledgeBase `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"knowledge_base,omitempty"`
	Messages      []ChatMessage  `gorm:"foreignKey:SessionID" json:"messages,omitempty"`
}

func CreateChatSession(session *ChatSession) error {
	return services.DB.Create(session).Error
}

func GetChatSessionByID(id uint) (*ChatSession, error) {
	var session ChatSession
	if err := services.DB.Preload("Messages").First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func ListChatSessionsByUserID(userID uint) ([]ChatSession, error) {
	var sessions []ChatSession
	if err := services.DB.Where("user_id = ?", userID).Order("updated_at DESC").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

func UpdateChatSession(session *ChatSession) error {
	return services.DB.Save(session).Error
}

func DeleteChatSession(id uint) error {
	return services.DB.Delete(&ChatSession{}, id).Error
}
