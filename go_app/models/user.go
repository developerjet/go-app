package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model        // 包含 ID、CreatedAt、UpdatedAt、DeletedAt
	Username     string `json:"username" gorm:"size:50;not null"`
	Email        string `json:"email" gorm:"size:100;not null;unique"`
	Password     string `json:"-" gorm:"size:255;not null"`  
	Token        string `json:"token,omitempty" gorm:"-"`    
	TokenVersion int    `json:"-" gorm:"default:0"`          
	AvatarURL    string `json:"avatarUrl" gorm:"column:avatar_url"` // 添加头像URL字段
}