package models

// PasswordChangeRequest 修改密码请求
type PasswordChangeRequest struct {
    OldPassword string `json:"oldPassword" binding:"required" example:"123456"`
    NewPassword string `json:"newPassword" binding:"required,min=6" example:"654321"`
}

// EmailUpdateRequest 更新邮箱请求
type EmailUpdateRequest struct {
    Email string `json:"email" binding:"required,email"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
    Username string `json:"username" binding:"required" example:"张三"`
    Email    string `json:"email" binding:"required,email" example:"zhangsan@example.com"`
    Password string `json:"password" binding:"required,min=6" example:"123456"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email" example:"zhangsan@example.com"`
    Password string `json:"password" binding:"required" example:"123456"`
}