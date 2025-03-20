package models

import "time"

type User struct {
    ID        uint      `json:"id" gorm:"primarykey"`
    Username  string    `json:"username"`
    Email     string    `json:"email" gorm:"type:varchar(255);unique"`
    Password  string    `json:"-"`  // 密码不返回
    Token     string    `json:"token,omitempty" gorm:"-"`  // token 不保存到数据库
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}