package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string     `gorm:"size:50;not null" json:"username"`
	Email     string     `gorm:"size:100;unique;not null" json:"email"`
	Password  string     `gorm:"size:255;not null" json:"-"`
	AvatarURL string     `gorm:"size:255" json:"avatar_url"`
	Birthday  *time.Time `gorm:"" json:"birthday"`         // 生日，允许为空
	Gender    string     `gorm:"size:10" json:"gender"`    // 性别：male/female/other
	Hobbies   string     `gorm:"size:500" json:"hobbies"`  // 爱好，用逗号分隔
	TokenVersion int        `gorm:"default:0" json:"-"`  // 添加 token 版本字段
	Token     string     `gorm:"-" json:"token,omitempty"` // 临时存储token
}
