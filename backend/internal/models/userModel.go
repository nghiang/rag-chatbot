package models

import (
	"time"
	"errors"
	"gorm.io/gorm"
	"backend/internal/database"
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

func MigrateUser() error {
	return database.DB.AutoMigrate(&User{})
}

func CreateUser(u *User) error {
	return database.DB.Create(u).Error
}

func GetUserByID(id uint) (*User, error) {
	var u User
	if err := database.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByEmail(email string) (*User, error) {
	var u User
	err := database.DB.Where("email = ?", email).First(&u).Error
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
	if err := database.DB.Select("id", "name").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
