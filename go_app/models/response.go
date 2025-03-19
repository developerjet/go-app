package models

// Response 通用响应结构
// @Description API 通用响应格式
type Response struct {
    Error   string      `json:"error,omitempty" example:"参数错误"`
    Message string      `json:"message,omitempty" example:"操作成功"`
    Data    interface{} `json:"data,omitempty" swaggertype:"object" example:"{}"`
}

// UserResponse 用户相关响应
// @Description 用户操作响应格式
type UserResponse struct {
    Message string `json:"message" example:"操作成功"`
    User    *User  `json:"user" swaggertype:"object"`
}

// LoginResponse 登录响应
// @Description 用户登录响应格式
type LoginResponse struct {
    Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
    UserInfo *User  `json:"userInfo" swaggertype:"object"`
}