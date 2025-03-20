package controllers

import (
	"go_app/models"
	"go_app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response "请求参数错误"
// @Failure 401 {object} models.Response "登录失败"
// @Router /login [post]
func (uc *UserController) Login(c *gin.Context) {
    var req models.LoginRequest
    contentType := c.GetHeader("Content-Type")
    
    var err error
    if contentType == "application/x-www-form-urlencoded" {
        err = c.ShouldBindWith(&req, binding.Form)
    } else {
        err = c.ShouldBindJSON(&req)
    }
    
    if err != nil {
        c.JSON(http.StatusOK, models.NewError("请求参数错误"))
        return
    }

    user, token, err := uc.userService.Login(req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusOK, models.NewError(err.Error()))
        return
    }

    loginResp := &models.LoginResponse{
        Token:    token,
        UserInfo: user,
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
    contentType := c.GetHeader("Content-Type")
    
    var err error
    if contentType == "application/x-www-form-urlencoded" {
        err = c.ShouldBindWith(&req, binding.Form)
    } else {
        err = c.ShouldBindJSON(&req)
    }
    
    if err != nil {
        c.JSON(http.StatusOK, models.NewError("请求参数错误"))
        return
    }

    // 检查邮箱是否已被注册
    if _, err := uc.userService.GetUserByEmail(req.Email); err == nil {
        c.JSON(http.StatusOK, models.NewError("该邮箱已被注册"))
        return
    }

    // 密码加密
    hashedPassword, err := uc.userService.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusOK, models.NewError("密码加密失败"))
        return
    }

    user := &models.User{
        Username: req.Username,
        Email:    req.Email,
        Password: hashedPassword,
    }

    if err := uc.userService.CreateUser(user); err != nil {
        c.JSON(http.StatusOK, models.NewError("创建用户失败"))
        return
    }
    c.JSON(http.StatusOK, models.NewSuccess(user, "注册成功"))
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
        ctx.JSON(http.StatusOK, models.NewError("请求参数错误"))
        return
    }

    if err := uc.userService.UpdateUserSafe(&user); err != nil {
        ctx.JSON(http.StatusOK, models.NewError("更新失败"))
        return
    }
    ctx.JSON(http.StatusOK, models.NewSuccess(user, "更新成功"))
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
        ctx.JSON(http.StatusOK, models.NewError("请求参数错误"))
        return
    }

    // 检查新邮箱是否已被使用
    if _, err := uc.userService.GetUserByEmail(req.Email); err == nil {
        ctx.JSON(http.StatusOK, models.NewError("该邮箱已被使用"))
        return
    }

    user, err := uc.userService.GetUserByID(req.UserID)
    if err != nil {
        ctx.JSON(http.StatusOK, models.NewError("用户不存在"))
        return
    }

    user.Email = req.Email
    if err := uc.userService.UpdateUserSafe(user); err != nil {
        ctx.JSON(http.StatusOK, models.NewError("更新邮箱失败"))
        return
    }

    ctx.JSON(http.StatusOK, models.NewSuccess(user, "邮箱更新成功"))
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
    contentType := ctx.GetHeader("Content-Type")
    
    var err error
    if contentType == "application/x-www-form-urlencoded" {
        err = ctx.ShouldBindWith(&req, binding.Form)
    } else {
        err = ctx.ShouldBindJSON(&req)
    }
    
    if err != nil {
        ctx.JSON(http.StatusOK, models.NewError("请求参数错误"))
        return
    }

    userID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusOK, models.NewError("未授权"))
        return
    }

    if err := uc.userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
        ctx.JSON(http.StatusOK, models.NewError(err.Error()))
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
        ctx.JSON(http.StatusOK, models.NewError("获取用户列表失败"))
        return
    }
    ctx.JSON(http.StatusOK, models.NewSuccess(users, "获取成功"))
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
        ctx.JSON(http.StatusOK, models.NewError("用户ID不能为空"))
        return
    }

    id, err := strconv.ParseUint(userID, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusOK, models.NewError("无效的用户ID"))
        return
    }

    user, err := uc.userService.GetUserByIDSafe(uint(id))
    if err != nil {
        ctx.JSON(http.StatusOK, models.NewError("用户不存在"))
        return
    }
    ctx.JSON(http.StatusOK, models.NewSuccess(user, "获取成功"))
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
        ctx.JSON(http.StatusOK, models.NewError("请求参数错误"))
        return
    }

    if err := uc.userService.DeleteUser(req.UserID); err != nil {
        ctx.JSON(http.StatusOK, models.NewError("删除用户失败"))
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
        ctx.JSON(http.StatusOK, models.NewError("未授权"))
        return
    }

    if err := uc.userService.Logout(userID.(uint)); err != nil {
        ctx.JSON(http.StatusOK, models.NewError("退出失败"))
        return
    }
    ctx.JSON(http.StatusOK, models.NewSuccess(nil, "退出成功"))
}
