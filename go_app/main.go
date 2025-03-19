// @title           Go App API
// @version         1.0
// @description     用户管理系统 API 文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"go_app/controllers"
	_ "go_app/docs" // 导入 swagger 文档
	"go_app/models"
	"go_app/routes"
	"go_app/services"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Go App API
// @version 1.0
// @description 用户管理系统 API 文档
// @host localhost:8080
// @BasePath /api
func main() {
	db, err := services.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	r := gin.Default()

	// 添加 JSON 相关中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
	})

	// Swagger 配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	api := r.Group("/api")
	{
		routes.SetupRoutes(api, userController) // 修改为正确的函数名
	}

	r.Run(":8080")
}
