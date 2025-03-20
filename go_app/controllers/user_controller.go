package controllers

import (
	"go_app/models"
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
// @Description 用户登录并获取token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body models.LoginRequest true "用户登录信息"
// @Success 200 {object} models.TokenResponse
// @Failure 401 {object} models.Response
// @Router /login [post]
func (uc *UserController) Login(c *gin.Context) {
    var loginRequest models.LoginRequest  // 使用 models 包中定义的类型，删除重复声明

    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest, models.Response{Error: "请求参数错误"})
        return
    }

    response, err := uc.userService.Login(loginRequest.Email, loginRequest.Password)
    // Login 方法中的响应格式
    if err != nil {
        c.JSON(http.StatusUnauthorized, models.Response{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

// ListUsers godoc
// @Summary 获取用户列表
// @Description 获取所有用户信息
// @Tags 用户管理
// @Produce json
// @Success 200 {object} models.Response{data=[]models.User}
// @Failure 500 {object} models.Response
// @Security ApiKeyAuth
// @Router /users [get]
// 将 'c' 改为 'uc' 以保持一致性
func (uc *UserController) ListUsers(ctx *gin.Context) {
    users, err := uc.userService.ListUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.Response{Data: users})
}

// GetUser godoc
// @Summary 获取单个用户
// @Description 根据ID获取用户信息
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} models.Response{data=models.User}
// @Failure 400,404 {object} models.Response
// @Security ApiKeyAuth
// @Router /users/{id} [get]
// GetUser godoc
func (uc *UserController) GetUser(ctx *gin.Context) {  // 将 c 改为 uc
    // 从路径参数获取用户ID
    userID := ctx.Param("id")
    if userID == "" {
        ctx.JSON(http.StatusBadRequest, models.Response{Error: "用户ID不能为空"})
        return
    }

    id, err := strconv.ParseUint(userID, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, models.Response{Error: "无效的用户ID"})
        return
    }

    user, err := uc.userService.GetUserByID(uint(id))  // 将 c 改为 uc
    if err != nil {
        ctx.JSON(http.StatusNotFound, models.Response{Error: "用户不存在"})
        return
    }

    ctx.JSON(http.StatusOK, models.Response{Data: user})
}

// UpdateUser godoc
// @Summary 更新用户信息
// @Description 更新用户基本信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body models.User true "用户信息"
// @Success 200 {object} models.Response{data=models.User}
// @Failure 400,500 {object} models.Response
// @Security ApiKeyAuth
// @Router /users/{id} [post]
// UpdateUser
func (uc *UserController) UpdateUser(ctx *gin.Context) {
    id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, models.Response{Error: "无效的用户ID"})
        return
    }
    var user models.User
    if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
        return
    }
    user.ID = uint(id)

    if err := uc.userService.UpdateUser(&user); err != nil {  // 将 c 改为 uc
        ctx.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, models.Response{Data: user})
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 根据ID删除用户
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Security ApiKeyAuth
// @Router /users/{id}/delete [post]
// DeleteUser
func (uc *UserController) DeleteUser(ctx *gin.Context) {
    id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, models.Response{Error: "无效的用户ID"})
        return
    }
    if err := uc.userService.DeleteUser(uint(id)); err != nil {  // 将 c 改为 uc
        ctx.JSON(http.StatusInternalServerError, models.Response{Error: err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, models.Response{Message: "用户已删除"})
}

// Register godoc
// @Summary 用户注册
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "用户注册信息"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.Response
// @Router /register [post]
func (uc *UserController) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.Response{Error: "请求参数错误"})
        return
    }

    // 检查邮箱是否已被注册
    if _, err := uc.userService.GetUserByEmail(req.Email); err == nil {
        c.JSON(http.StatusBadRequest, models.Response{
            Error: "该邮箱已被注册",
        })
        return
    }

    // 密码加密
    hashedPassword, err := uc.userService.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.Response{
            Error: "密码加密失败",
        })
        return
    }

    user := &models.User{
        Username: req.Username,
        Email:    req.Email,
        Password: hashedPassword,
    }

    if err := uc.userService.CreateUser(user); err != nil {
        c.JSON(http.StatusInternalServerError, models.Response{
            Error: "创建用户失败",
        })
        return
    }

    c.JSON(http.StatusOK, models.Response{
        Message: "注册成功",
        Data:    user,
    })
}

// UpdateEmail godoc
// @Summary 更新用户邮箱
// @Description 更新指定用户的邮箱地址
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param email body models.EmailUpdate true "新邮箱"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.Response
// @Security ApiKeyAuth
// @Router /users/{id}/email [post]
// UpdateEmail
func (uc *UserController) UpdateEmail(ctx *gin.Context) {
    // 获取用户ID
    userID := ctx.Param("id")
    if userID == "" {
        ctx.JSON(http.StatusBadRequest, models.Response{Error: "用户ID不能为空"})
        return
    }

    id, err := strconv.ParseUint(userID, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, models.Response{Error: "无效的用户ID"})
        return
    }

    // 获取新邮箱
    // 删除本地 EmailUpdate 结构体定义
    var update models.EmailUpdate  // 直接使用 models 包中定义的类型
    if err := ctx.ShouldBindJSON(&update); err != nil {
        ctx.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
        return
    }

    // 检查新邮箱是否已被使用
    if _, err := uc.userService.GetUserByEmail(update.Email); err == nil {  // 将 c 改为 uc
        ctx.JSON(http.StatusBadRequest, models.Response{Error: "该邮箱已被使用"})
        return
    }

    user, err := uc.userService.GetUserByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, models.Response{Error: "用户不存在"})
        return
    }

    if err := uc.userService.UpdateUser(user); err != nil {  // 将 c 改为 uc
        ctx.JSON(http.StatusInternalServerError, models.Response{Error: "更新邮箱失败"})
        return
    }

    ctx.JSON(http.StatusOK, models.UserResponse{
        Message: "邮箱更新成功",
        User:    user,
    })
}

// ChangePassword godoc
// @Summary 修改用户密码
// @Description 用户修改自己的密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param Authorization header string true "用户令牌"
// @Param request body models.PasswordChangeRequest true "密码修改请求"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "请求参数错误"
// @Failure 401 {object} models.Response "未授权"
// @Router /api/users/password [post]
func (uc *UserController) ChangePassword(ctx *gin.Context) {
    var req models.PasswordChangeRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, models.Response{Error: "请提供正确的密码格式"})
        return
    }

    // 从上下文中获取用户ID
    userID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, models.Response{Error: "未授权"})
        return
    }

    // 修改类型转换方式
    id, ok := userID.(uint)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, models.Response{Error: "用户ID类型错误"})
        return
    }

    // 使用正确的类型获取用户信息
    user, err := uc.userService.GetUserByID(id)
    if err != nil {
        ctx.JSON(http.StatusNotFound, models.Response{Error: "用户不存在"})
        return
    }

    // 验证旧密码是否正确
    if err := uc.userService.VerifyPassword(user.Password, req.OldPassword); err != nil {
        ctx.JSON(http.StatusBadRequest, models.Response{Error: "旧密码错误"})
        return
    }
    
    err = uc.userService.ChangePassword(user.ID, req.OldPassword, req.NewPassword)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, models.Response{Error: err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, models.Response{Message: "密码修改成功"})
}

// Logout godoc
// @Summary 用户退出
// @Description 用户退出登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.Response "退出成功"
// @Failure 401 {object} models.Response "未授权"
// @Failure 500 {object} models.Response "服务器错误"
// @Router /api/users/logout [post]  // 修改路由路径，添加 /api 前缀
func (uc *UserController) Logout(ctx *gin.Context) {
    // 从上下文获取用户ID
    userID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, models.Response{
            Error: "未授权",
        })
        return
    }

    // 调用服务层的退出方法
    if err := uc.userService.Logout(userID.(uint)); err != nil {
        ctx.JSON(http.StatusInternalServerError, models.Response{
            Error: "退出失败",
        })
        return
    }

    ctx.JSON(http.StatusOK, models.Response{
        Message: "退出成功",
    })
}
