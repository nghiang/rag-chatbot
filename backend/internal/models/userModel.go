package models

import (
	"time"
	"errors"
	"gorm.io/gorm"
	"backend/internal/services"
)

type User struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    Name         string    `json:"name"`
    Email        string    `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string    `gorm:"size:255;not null" json:"-"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`

    KnowledgeBases []KnowledgeBase `gorm:"foreignKey:UserID" json:"knowledge_bases,omitempty"`
    ChatSessions   []ChatSession     `gorm:"foreignKey:UserID" json:"chat_sessions,omitempty"`
}

func CreateUser(u *User) error {
	return services.DB.Create(u).Error
}

func GetUserByID(id uint) (*User, error) {
	var u User
	if err := services.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByEmail(email string) (*User, error) {
	var u User
	err := services.DB.Where("email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func ListUsers() ([]User, error) {
    var users []User
    if err := services.DB.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}
