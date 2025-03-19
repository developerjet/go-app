package models

import (
    "time"
)

// User 用户模型
// @Description 用户信息
type User struct {
    ID        uint      `json:"id" gorm:"primarykey"`
    Username  string    `json:"username" gorm:"column:username;not null"`
    Email     string    `json:"email" gorm:"column:email;unique;not null"`
    Password  string    `json:"-" gorm:"column:password;not null"`
    CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
    UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

// EmailUpdate 邮箱更新请求
// @Description 用户邮箱更新请求参数
type EmailUpdate struct {
    Email string `json:"email" binding:"required,email" example:"new_email@example.com"`
}