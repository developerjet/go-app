package routes

import (
	"go_app/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.RouterGroup, userController *controllers.UserController) {
	// 用户相关路由
	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
	
	// 需要认证的路由
	r.GET("/users", userController.ListUsers)
	r.GET("/users/:id", userController.GetUser)
	r.POST("/users/:id", userController.UpdateUser)
	r.POST("/users/:id/delete", userController.DeleteUser)
	r.POST("/users/:id/email", userController.UpdateEmail)
}
