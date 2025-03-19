package routes

import (
	"go_app/controllers"
	"go_app/middleware"  // 添加 middleware 包的导入

	"github.com/gin-gonic/gin"
)

func SetupRoutes(api *gin.RouterGroup, userController *controllers.UserController) {
    // 无需认证的路由
    api.POST("/register", userController.Register)
    api.POST("/login", userController.Login)

    // 需要认证的路由
    users := api.Group("/users")
    users.Use(middleware.AuthMiddleware()) // 添加认证中间件
    {
        users.GET("", userController.ListUsers)            // 获取用户列表
        users.GET("/:id", userController.GetUser)         // 获取单个用户
        users.POST("/:id", userController.UpdateUser)     // 更新用户
        users.POST("/:id/delete", userController.DeleteUser) // 删除用户
        users.POST("/:id/email", userController.UpdateEmail) // 更新邮箱
        users.POST("/password", userController.ChangePassword)
    }
}
