package controllers

import (
	"go_app/models"
	"go_app/pkg/errcode"
	"go_app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录并获取token，同一用户重复登录会使之前的token失效
// @Tags 用户管理
// @Accept json,x-www-form-urlencoded
// @Produce json
// @Param user body models.LoginRequest true "用户登录信息"
// @Success 200 {object} models.Response{data} "登录成功"
// @Failure 400 {object} models.Response "请求参数错误"
// @Failure 1000 {object} models.Response "用户不存在"
// @Failure 1001 {object} models.Response "密码错误"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /login [post]
func (uc *UserController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, models.NewError(errcode.InvalidParams))
		return
	}

	user, token, err := uc.userService.Login(req.Email, req.Password)
	if err != nil {
		if err.Error() == "用户不存在" {
			c.JSON(http.StatusOK, models.NewError(errcode.UserNotFound))
			return
		}
		if err.Error() == "密码错误" {
			c.JSON(http.StatusOK, models.NewError(errcode.UserPasswordError))
			return
		}
		c.JSON(http.StatusOK, models.NewError(errcode.ServerError))
		return
	}

	// 将 token 只设置到 user 对象中
	user.Token = token

	// 在 Login 方法中
	loginResp := &models.LoginResponse{
	    UserInfo: user.ToUserInfo(),
	}

	c.JSON(http.StatusOK, models.NewSuccess(loginResp, "登录成功"))
}

// Register godoc
// @Summary 用户注册
// @Description 新用户注册
// @Tags 用户管理
// @Accept json,x-www-form-urlencoded
// @Produce json
// @Param user body models.RegisterRequest true "用户注册信息"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "请求参数错误"
// @Failure 409 {object} models.Response "邮箱已被注册"
// @Router /register [post]
func (uc *UserController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, models.NewError(errcode.InvalidParams))
		return
	}

	if _, err := uc.userService.GetUserByEmail(req.Email); err == nil {
		c.JSON(http.StatusOK, models.NewError(errcode.UserAlreadyExists))
		return
	}

	hashedPassword, err := uc.userService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, models.NewError(errcode.ServerError))
		return
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := uc.userService.CreateUser(user); err != nil {
		c.JSON(http.StatusOK, models.NewError(errcode.UserCreateFailed))
		return
	}
	c.JSON(http.StatusOK, models.NewSuccess(user.ToUserInfo(), "注册成功"))
}

// UpdateUser godoc
// @Summary 更新用户信息
// @Description 更新用户基本信息
// @Tags 用户管理
// @Accept json,x-www-form-urlencoded
// @Produce json
// @Param user body models.User true "用户信息"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "请求参数错误"
// @Failure 404 {object} models.Response "用户不存在"
// @Security ApiKeyAuth
// @Router /users/update [post]
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.InvalidParams))
		return
	}

	if err := uc.userService.UpdateUserSafe(&user); err != nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.UserUpdateFailed))
		return
	}
	ctx.JSON(http.StatusOK, models.NewSuccess(user.ToUserInfo(), "更新成功"))
}

// UpdateEmail godoc
// @Summary 更新邮箱
// @Description 更新用户邮箱
// @Tags 用户管理
// @Accept json,x-www-form-urlencoded
// @Produce json
// @Param user body models.EmailUpdateRequest true "更新邮箱信息"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "请求参数错误"
// @Failure 409 {object} models.Response "邮箱已被使用"
// @Security ApiKeyAuth
// @Router /users/email [post]
func (uc *UserController) UpdateEmail(ctx *gin.Context) {
	var req models.EmailUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.InvalidParams))
		return
	}

	if _, err := uc.userService.GetUserByEmail(req.Email); err == nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.UserAlreadyExists))
		return
	}

	user, err := uc.userService.GetUserByID(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.UserNotFound))
		return
	}

	user.Email = req.Email
	if err := uc.userService.UpdateUserSafe(user); err != nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.UserUpdateFailed))
		return
	}

	ctx.JSON(http.StatusOK, models.NewSuccess(user.ToUserInfo(), "邮箱更新成功"))
}

// ChangePassword godoc
// @Summary 修改密码
// @Description 用户修改密码
// @Tags 用户管理
// @Accept json,x-www-form-urlencoded
// @Produce json
// @Param user body models.PasswordChangeRequest true "修改密码信息"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "请求参数错误"
// @Failure 401 {object} models.Response "未授权"
// @Security ApiKeyAuth
// @Router /users/password [post]
func (uc *UserController) ChangePassword(ctx *gin.Context) {
	var req models.PasswordChangeRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.InvalidParams))
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, models.NewError(errcode.Unauthorized))
		return
	}

	if err := uc.userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		// 这里错误地使用了 c 而不是 ctx
		ctx.JSON(http.StatusOK, models.NewError(errcode.UserPasswordError))
		return
	}
	ctx.JSON(http.StatusOK, models.NewSuccess(nil, "密码修改成功"))
}

// ListUsers godoc
// @Summary 获取用户列表
// @Description 获取所有用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "未授权"
// @Security ApiKeyAuth
// @Router /users [get]
func (uc *UserController) ListUsers(ctx *gin.Context) {
    users, err := uc.userService.ListUsers()
    if err != nil {
        ctx.JSON(http.StatusOK, models.NewError(errcode.ServerError))
        return
    }
    
    // 转换用户列表
    userInfos := make([]*models.UserInfo, len(users))
    for i, user := range users {
        userInfos[i] = user.ToUserInfo()
    }
    
    ctx.JSON(http.StatusOK, models.NewSuccess(userInfos, "获取成功"))
}

// GetUser godoc
// @Summary 获取用户信息
// @Description 获取指定用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "请求参数错误"
// @Failure 404 {object} models.Response "用户不存在"
// @Security ApiKeyAuth
// @Router /users/info/{id} [get]
func (uc *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusOK, models.NewError(errcode.InvalidParams))
		return
	}

	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.InvalidParams))
		return
	}

	user, err := uc.userService.GetUserByIDSafe(uint(id))
	if err != nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.UserNotFound))
		return
	}
	// GetUser 方法中
	ctx.JSON(http.StatusOK, models.NewSuccess(user.ToUserInfo(), "获取成功"))
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 删除指定用户
// @Tags 用户管理
// @Accept json,x-www-form-urlencoded
// @Produce json
// @Param user body models.UserIDRequest true "用户ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "请求参数错误"
// @Failure 404 {object} models.Response "用户不存在"
// @Security ApiKeyAuth
// @Router /users/delete [post]
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	var req models.UserIDRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.InvalidParams))
		return
	}

	if err := uc.userService.DeleteUser(req.UserID); err != nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.UserDeleteFailed))
		return
	}
	ctx.JSON(http.StatusOK, models.NewSuccess(nil, "删除成功"))
}

// Logout godoc
// @Summary 用户退出
// @Description 用户退出登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response "未授权"
// @Security ApiKeyAuth
// @Router /logout [post]
func (uc *UserController) Logout(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusOK, models.NewError(errcode.Unauthorized))
		return
	}

	if err := uc.userService.Logout(userID.(uint)); err != nil {
		ctx.JSON(http.StatusOK, models.NewError(errcode.ServerError))
		return
	}
	ctx.JSON(http.StatusOK, models.NewSuccess(nil, "退出成功"))
}
