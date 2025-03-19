package models

// Response 通用响应结构
type Response struct {
    Error   string      `json:"error,omitempty" example:"参数错误"`
    Message string      `json:"message,omitempty" example:"操作成功"`
    Data    interface{} `json:"data,omitempty"`
}

// UserResponse 用户相关响应
type UserResponse struct {
    Message string `json:"message" example:"操作成功"`
    User    *User  `json:"user" swaggertype:"object"`
}

// LoginResponse 登录响应
type LoginResponse struct {
    Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
    UserInfo *User  `json:"userInfo" swaggertype:"object"`
}

// TokenResponse 登录成功返回的 token 响应
type TokenResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}