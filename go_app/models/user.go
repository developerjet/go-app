package models

// User 用户模型
// @Description 用户信息
type User struct {
    ID       uint   `json:"id" gorm:"primarykey" example:"1"`
    Name     string `json:"name" example:"张三"`
    Email    string `json:"email" gorm:"unique" example:"zhangsan@example.com"`
    Password string `json:"-"` // 密码不返回给前端
    Age      int    `json:"age" example:"25"`
}

// LoginRequest 登录请求
// @Description 用户登录请求参数
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email" example:"zhangsan@example.com"`
    Password string `json:"password" binding:"required" example:"123456"`
}

// RegisterRequest 注册请求
// @Description 用户注册请求参数
type RegisterRequest struct {
    Name     string `json:"name" binding:"required" example:"张三"`
    Email    string `json:"email" binding:"required,email" example:"zhangsan@example.com"`
    Password string `json:"password" binding:"required" example:"123456"`
    Age      int    `json:"age" binding:"required" example:"25"`
}

// EmailUpdate 邮箱更新请求
type EmailUpdate struct {
    Email string `json:"email" binding:"required,email" example:"new_email@example.com"`
}