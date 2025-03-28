package models

import (
    "go_app/pkg/errcode"
    "time"
)

// Response 统一响应结构
// @Description API 统一响应格式
type Response struct {
    Code      int         `json:"code" example:"200" description:"状态码 200成功，其他失败"`
    Message   string      `json:"message" example:"操作成功" description:"提示信息"`
    Data      interface{} `json:"data,omitempty" description:"响应数据"`
    Timestamp int64      `json:"timestamp" example:"1704067200" description:"时间戳"`
}

// TokenResponse 登录成功返回的 token 响应
// @Description 登录成功返回的访问令牌
type TokenResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." description:"JWT访问令牌"`
}

// UserInfo 用户信息响应结构体
// @Description 用户详细信息响应结构
type UserInfo struct {
    UserID    uint       `json:"userId" example:"1" description:"用户ID"`
    Username  string     `json:"username" example:"张三" description:"用户名"`
    Email     string     `json:"email" example:"zhangsan@example.com" description:"邮箱地址"`
    AvatarURL string     `json:"avatarUrl" example:"https://example.com/avatar.jpg" description:"头像URL"`
    Birthday  *time.Time `json:"birthday" example:"1990-01-01T00:00:00+08:00" description:"生日"`
    Gender    string     `json:"gender" example:"male" enums:"male,female,other" description:"性别"`
    Hobbies   string     `json:"hobbies" example:"读书,游泳,旅行" description:"用户爱好，多个爱好用逗号分隔"`
    CreatedAt time.Time  `json:"createdAt" example:"2024-01-01T00:00:00+08:00" description:"创建时间"`
    UpdatedAt time.Time  `json:"updatedAt" example:"2024-01-01T00:00:00+08:00" description:"更新时间"`
}

// LoginResponse 登录响应
// @Description 用户登录成功响应
type LoginResponse struct {
    UserInfo *UserInfo `json:"userInfo" description:"用户信息"`
}

// Pagination 分页信息
// @Description 分页查询信息
type Pagination struct {
    Page     int   `json:"page" example:"1" minimum:"1" description:"当前页码"`
    PageSize int   `json:"pageSize" example:"10" minimum:"1" maximum:"100" description:"每页数量"`
    Total    int64 `json:"total" example:"100" description:"总记录数"`
}

// PageResponse 分页响应
// @Description 通用分页响应结构
type PageResponse struct {
    List       interface{} `json:"list" description:"数据列表"`
    Pagination Pagination  `json:"pagination" description:"分页信息"`
}

// UserPageResponse 用户分页列表响应
// @Description 用户分页列表响应结构
type UserPageResponse struct {
    List     []*UserInfo `json:"list" description:"用户列表数据"`
    Page     int         `json:"page" example:"1" minimum:"1" description:"当前页码"`
    PageSize int         `json:"pageSize" example:"10" minimum:"1" maximum:"100" description:"每页数量"`
    Total    int64       `json:"total" example:"100" description:"总记录数"`
}

// TokenStatusResponse Token状态响应
// @Description Token状态信息
type TokenStatusResponse struct {
    Status  int    `json:"status" example:"4001" description:"token状态码"`
    Message string `json:"message" example:"无效的访问凭证" description:"状态描述"`
}

// Token 状态码
const (
    TokenStatusValid   = 1    // token有效
    TokenStatusInvalid = 4001 // token无效
    TokenStatusExpired = 4002 // token过期
    TokenStatusReplace = 4003 // token被替换（在其他设备登录）
)

// TokenErrorMessages token错误信息映射
var TokenErrorMessages = map[int]string{
    TokenStatusInvalid: "无效的访问凭证",
    TokenStatusExpired: "访问凭证已过期",
    TokenStatusReplace: "您的账号已在其他设备登录",
}

// 响应构造函数
func NewSuccess(data interface{}, message string) Response {
    return Response{
        Code:      200,
        Message:   message,
        Data:      data,
        Timestamp: time.Now().Unix(),
    }
}

func NewError(err *errcode.ErrorCode) *Response {
    return &Response{
        Code:      err.Code,
        Message:   err.Message,
        Timestamp: time.Now().Unix(),
    }
}

func NewTokenError(status int) *Response {
    message := TokenErrorMessages[status]
    return &Response{
        Code:    status,
        Message: message,
        Data: TokenStatusResponse{
            Status:  status,
            Message: message,
        },
        Timestamp: time.Now().Unix(),
    }
}

// WithDetails 为响应添加详细信息
func (r *Response) WithDetails(details string) *Response {
    r.Message = r.Message + ": " + details
    return r
}

// ToUserInfo 用户模型转换为用户信息
func (u *User) ToUserInfo() *UserInfo {
    return &UserInfo{
        UserID:    u.ID,
        Username:  u.Username,
        Email:     u.Email,
        AvatarURL: u.AvatarURL,
        Birthday:  u.Birthday,
        Gender:    u.Gender,
        Hobbies:   u.Hobbies,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}
