package models

import (
	"backend/internal/services"
	"time"
)

type ChatMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"index;not null" json:"session_id"`
	Role      string    `gorm:"size:20;not null" json:"role"` // "user" | "assistant" | "system"
	Message   string    `gorm:"type:text" json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateChatMessage(message *ChatMessage) error {
	return services.DB.Create(message).Error
}

func GetChatMessageByID(id uint) (*ChatMessage, error) {
	var message ChatMessage
	if err := services.DB.First(&message, id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func ListChatMessagesBySessionID(sessionID uint) ([]ChatMessage, error) {
	var messages []ChatMessage
	if err := services.DB.Where("session_id = ?", sessionID).Order("created_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func UpdateChatMessage(message *ChatMessage) error {
	return services.DB.Save(message).Error
}

func DeleteChatMessage(id uint) error {
	return services.DB.Delete(&ChatMessage{}, id).Error
}
