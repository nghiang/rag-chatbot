package models

import (
	"time"
	// "errors"
	// "gorm.io/gorm"
	"backend/internal/services"
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

func CreateKnowledgeBase(kb *KnowledgeBase) error {
    return services.DB.Create(kb).Error
}

func GetKnowledgeBaseByID(id uint) (*KnowledgeBase, error) {
    var kb KnowledgeBase
    if err := services.DB.First(&kb, id).Error; err != nil {
        return nil, err
    }
    return &kb, nil
}

func ListKnowledgeBasesByUserID(userID uint) ([]KnowledgeBase, error) {
    var kbs []KnowledgeBase
    if err := services.DB.Where("user_id = ?", userID).Find(&kbs).Error; err != nil {
        return nil, err
    }
    return kbs, nil
}

func UpdateKnowledgeBase(kb *KnowledgeBase) error {
    return services.DB.Save(kb).Error
}

func DeleteKnowledgeBase(id uint) error {
    return services.DB.Delete(&KnowledgeBase{}, id).Error
}


