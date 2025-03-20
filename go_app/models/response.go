package models

import "time"

// Response 统一响应结构
type Response struct {
    Code      int         `json:"code"`              // 状态码 200成功，其他失败
    Message   string      `json:"message"`           // 提示信息
    Data      interface{} `json:"data,omitempty"`    // 数据
    Timestamp int64       `json:"timestamp"`         // 时间戳
}

// NewSuccess 成功响应
func NewSuccess(data interface{}, message string) Response {
    return Response{
        Code:      200,
        Message:   message,
        Data:      data,
        Timestamp: time.Now().Unix(),
    }
}

// NewError 错误响应
func NewError(message string) Response {
    return Response{
        Code:      500,
        Message:   message,
        Timestamp: time.Now().Unix(),
    }
}

// UserResponse 用户相关响应
type UserResponse struct {
    Message string `json:"message" example:"操作成功"`
    User    *User  `json:"user" swaggertype:"object"`
}

// Token 状态码
const (
    TokenStatusValid   = 1    // token有效
    TokenStatusInvalid = 4001 // token无效
    TokenStatusExpired = 4002 // token过期
    TokenStatusReplace = 4003 // token被替换（在其他设备登录）
)

// Token 错误提示信息
var TokenErrorMessages = map[int]string{
    TokenStatusInvalid: "无效的访问凭证",
    TokenStatusExpired: "访问凭证已过期",
    TokenStatusReplace: "您的账号已在其他设备登录",
}

// TokenStatusResponse Token状态响应
type TokenStatusResponse struct {
    Status  int    `json:"status"`   // token状态码
    Message string `json:"message"`   // 状态描述
}

// NewTokenError 创建Token错误响应
func NewTokenError(status int) Response {
    message := TokenErrorMessages[status]
    return Response{
        Code:      status,
        Message:   message,
        Data:      TokenStatusResponse{
            Status:  status,
            Message: message,
        },
        Timestamp: time.Now().Unix(),
    }
}

// TokenResponse 登录成功返回的 token 响应
type TokenResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// LoginResponse 登录响应
type LoginResponse struct {
    Token    string `json:"token"`
    UserInfo *User  `json:"userInfo"`
}