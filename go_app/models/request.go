package models

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

// UserIDRequest 用户ID请求参数
type UserIDRequest struct {
    UserID uint `json:"userId" form:"userId" binding:"required"`
}

// EmailUpdateRequest 更新邮箱请求
type EmailUpdateRequest struct {
    UserID uint   `json:"userId" form:"userId" binding:"required"`
    Email  string `json:"email" form:"email" binding:"required,email"`
}