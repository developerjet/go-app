package errcode

// ErrorCode 错误码结构体
type ErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 系统级错误码
var (
	Success         = &ErrorCode{Code: 200, Message: "成功"}
	ServerError     = &ErrorCode{Code: 500, Message: "服务器内部错误"}
	InvalidParams   = &ErrorCode{Code: 400, Message: "请求参数错误"}
	NotFound        = &ErrorCode{Code: 404, Message: "资源不存在"}
	Unauthorized    = &ErrorCode{Code: 401, Message: "未授权"}
	TooManyRequests = &ErrorCode{Code: 429, Message: "请求过于频繁"}

	// 用户模块错误码 (1000-1999)
	UserNotFound      = &ErrorCode{Code: 1000, Message: "用户不存在"}
	UserPasswordError = &ErrorCode{Code: 1001, Message: "密码错误"}
	UserAlreadyExists = &ErrorCode{Code: 1002, Message: "用户已存在"}
	UserCreateFailed  = &ErrorCode{Code: 1003, Message: "创建用户失败"}
	UserUpdateFailed  = &ErrorCode{Code: 1004, Message: "更新用户信息失败"}
	UserDeleteFailed  = &ErrorCode{Code: 1005, Message: "删除用户失败"}

	// Token相关错误码 (2000-2999)
	TokenInvalid      = &ErrorCode{Code: 2000, Message: "无效的Token"}
	TokenExpired      = &ErrorCode{Code: 2001, Message: "Token已过期"}
	TokenMissing      = &ErrorCode{Code: 2002, Message: "请提供认证Token"}
	TokenVersionError = &ErrorCode{Code: 2003, Message: "Token已失效，请重新登录"}

	InvalidRequest = &ErrorCode{Code: 40001, Message: "无效的请求"}
)

// Error 实现error接口
func (e *ErrorCode) Error() string {
	return e.Message
}

// NewError 创建新的错误响应
func NewError(code int, message string) *ErrorCode {
	return &ErrorCode{
		Code:    code,
		Message: message,
	}
}
