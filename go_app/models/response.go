package models

import (
	"go_app/pkg/errcode"
	"time"
)

// Response 统一响应结构
type Response struct {
	Code      int         `json:"code"`           // 状态码 200成功，其他失败
	Message   string      `json:"message"`        // 提示信息
	Data      interface{} `json:"data,omitempty"` // 数据
	Timestamp int64       `json:"timestamp"`      // 时间戳
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

// NewError 使用 ErrorCode 创建错误响应
func NewError(err *errcode.ErrorCode) Response {
	return Response{
		Code:      err.Code,
		Message:   err.Message,
		Timestamp: time.Now().Unix(),
	}
}

// UserResponse 用户相关响应
type UserResponse struct {
	Message string    `json:"message" example:"操作成功"`
	User    *UserInfo `json:"user" swaggertype:"object"`
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
	Status  int    `json:"status"`  // token状态码
	Message string `json:"message"` // 状态描述
}

// NewTokenError 创建Token错误响应
func NewTokenError(status int) Response {
	message := TokenErrorMessages[status]
	return Response{
		Code:    status,
		Message: message,
		Data: TokenStatusResponse{
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

// UserInfo 用户信息响应结构体
type UserInfo struct {
    UserId    uint      `json:"userId"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
    DeletedAt string    `json:"deletedAt,omitempty"` // 改为 string 类型
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Token     string    `json:"token,omitempty"`
}

// 添加转换方法
func (u *User) ToUserInfo() *UserInfo {
    info := &UserInfo{
        UserId:    u.ID,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
        Username:  u.Username,
        Email:     u.Email,
        Token:     u.Token,
    }
    
    // 处理软删除时间
    if u.DeletedAt.Valid {
        info.DeletedAt = u.DeletedAt.Time.Format(time.RFC3339)
    }
    
    return info
}

// LoginResponse 登录响应
type LoginResponse struct {
    UserInfo *UserInfo `json:"userInfo"`
}


// Pagination 分页信息
type Pagination struct {
    Page      int   `json:"page"`      // 当前页码
    PageSize  int   `json:"pageSize"`  // 每页数量
    Total     int64 `json:"total"`     // 总记录数
}

// PageResponse 分页响应
type PageResponse struct {
    List       interface{} `json:"list"`       // 数据列表
    Pagination Pagination  `json:"pagination"` // 分页信息
}

// UserPageResponse 用户分页列表响应
type UserPageResponse struct {
    List     []*UserInfo `json:"list"`      // 用户列表
    Page     int         `json:"page"`      // 当前页码
    PageSize int         `json:"pageSize"`  // 每页数量
    Total    int64       `json:"total"`     // 总记录数
}

// 删除不需要的结构体
// 删除 Pagination 结构体
// 删除 PageResponse 结构体
