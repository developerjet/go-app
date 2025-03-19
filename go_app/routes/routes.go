package routes

import (
	"go_app/controllers"
	"go_app/middleware"  // 添加 middleware 包的导入

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置所有的路由
// @param api *gin.RouterGroup - gin路由组指针，用于注册路由
// @param userController *controllers.UserController - 用户控制器指针，处理具体的业务逻辑
func SetupRoutes(api *gin.RouterGroup, userController *controllers.UserController) {
    // 无需认证的路由组
    // POST /register - 用户注册接口
    // 接收用户名、邮箱和密码，创建新用户
    api.POST("/register", userController.Register)

    // POST /login - 用户登录接口
    // 接收邮箱和密码，返回JWT token
    api.POST("/login", userController.Login)

    // 需要认证的路由组
    users := api.Group("/users")
    users.Use(middleware.AuthMiddleware()) // 添加JWT认证中间件，验证请求头中的token
    {
        // GET /users - 获取所有用户列表
        // 需要管理员权限
        users.GET("", userController.ListUsers)

        // GET /users/:id - 获取指定ID的用户信息
        // URL参数: id - 用户ID
        users.GET("/:id", userController.GetUser)

        // POST /users/:id - 更新指定用户的基本信息
        // URL参数: id - 用户ID
        // Body: 用户信息（用户名等）
        users.POST("/:id", userController.UpdateUser)

        // POST /users/:id/delete - 删除指定用户
        // URL参数: id - 用户ID
        users.POST("/:id/delete", userController.DeleteUser)

        // POST /users/:id/email - 更新用户邮箱
        // URL参数: id - 用户ID
        // Body: 新的邮箱地址
        users.POST("/:id/email", userController.UpdateEmail)

        // POST /users/password - 修改当前用户密码
        // Body: 旧密码和新密码
        // 注意：此接口使用token中的用户ID，不需要在URL中指定用户ID
        users.POST("/password", userController.ChangePassword)
    }
}
