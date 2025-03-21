package models

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model        // 这个已经包含了 ID、CreatedAt、UpdatedAt、DeletedAt
    Username     string `json:"username" gorm:"size:50;not null"`
    Email        string `json:"email" gorm:"size:100;not null;unique"`
    Password     string `json:"-" gorm:"size:255;not null"`  // 密码不返回
    Token        string `json:"token,omitempty" gorm:"-"`    // token 不保存到数据库
    TokenVersion int    `json:"-" gorm:"default:0"`          // 用于追踪 token 版本
}