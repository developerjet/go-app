package models

import "time"

// LoginRequest 登录请求参数
type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email" example:"zhangsan@example.com"`
	Password string `json:"password" form:"password" binding:"required" example:"123456"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" form:"username" binding:"required" example:"张三"`
	Email    string `json:"email" form:"email" binding:"required,email" example:"zhangsan@example.com"`
	Password string `json:"password" form:"password" binding:"required,min=6" example:"123456"`
}

// PasswordChangeRequest 修改密码请求
type PasswordChangeRequest struct {
	OldPassword string `json:"oldPassword" form:"oldPassword" binding:"required" example:"123456"`
	NewPassword string `json:"newPassword" form:"newPassword" binding:"required,min=6" example:"654321"`
}

// UserIDRequest 用户ID请求
type UserIDRequest struct {
    UserID uint `json:"userId" binding:"required" example:"1" description:"用户ID"`
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
    UserID    uint       `json:"userId" binding:"required" example:"1" description:"用户ID"`
    Username  string     `json:"username" binding:"required,min=2,max=50" example:"张三" description:"用户名"`
    Birthday  *time.Time `json:"birthday" example:"1990-01-01T00:00:00+08:00" description:"生日"`
    Gender    string     `json:"gender" enums:"male,female,other" example:"male" description:"性别"`
    Hobbies   string     `json:"hobbies" example:"读书,游泳,旅行" description:"爱好，多个爱好用逗号分隔"`
}

// EmailUpdateRequest 更新邮箱请求
type EmailUpdateRequest struct {
    UserID uint   `json:"userId" binding:"required" example:"1" description:"用户ID"`
    Email  string `json:"email" binding:"required,email" example:"newemail@example.com" description:"新邮箱地址"`
}
