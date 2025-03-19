package models

// Response 通用响应结构
type Response struct {
    Message string      `json:"message,omitempty" example:"操作成功"`
    Error   string      `json:"error,omitempty" example:"错误信息"`
    Data    interface{} `json:"data,omitempty"`
}

// TokenResponse 登录响应
type TokenResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UserResponse 用户响应
type UserResponse struct {
    Message string `json:"message" example:"操作成功"`
    User    *User  `json:"user"`
}