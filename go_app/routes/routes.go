package routes

import (
	"go_app/controllers"
	"go_app/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
// @Summary 设置API路由
// @Description 配置所有API路由和中间件
// @Tags 系统
// @Accept json
// @Produce json
// @Param userController body controllers.UserController true "用户控制器"
// @Router /api [post]
func SetupRoutes(api *gin.RouterGroup, userController *controllers.UserController) {
    // 无需认证的路由组
    api.POST("/register", userController.Register)
    api.POST("/login", userController.Login)

    // 需要认证的路由组
    users := api.Group("/users")
    users.Use(middleware.AuthMiddleware())
    {
        users.GET("", userController.ListUsers)
        users.GET("/info", userController.GetUser)         // 改为 query 参数
        users.POST("/update", userController.UpdateUser)   // 改为 body 参数
        users.POST("/delete", userController.DeleteUser)   // 改为 body 参数
        users.POST("/email", userController.UpdateEmail)   // 改为 body 参数
        users.POST("/password", userController.ChangePassword)
        users.POST("/logout", userController.Logout)
    }
}
