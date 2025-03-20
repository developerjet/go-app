// @title Go App API
// @version 1.0
// @description 用户管理系统 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
    "fmt"
    "go_app/config"
    "go_app/controllers"
    _ "go_app/docs"
    "go_app/models"
    "go_app/routes"
    "go_app/services"
    "log"

    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
    // 加载配置
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("加载配置失败:", err)
    }

    // 设置 Gin 模式
    if cfg.Server.Mode == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    // 连接数据库
    if err := services.ConnectDB(); err != nil {
        log.Fatal("连接数据库失败:", err)
    }

    // 获取数据库实例
    db := services.GetDB()

    // 自动迁移数据库表结构
    if err := db.AutoMigrate(&models.User{}, &models.UserToken{}); err != nil {
        log.Fatal("数据库迁移失败:", err)
    }

    // 初始化服务和控制器
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

    // 使用配置的端口启动服务器
    serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
    log.Printf("服务器启动在 http://localhost%s", serverAddr)
    if err := r.Run(serverAddr); err != nil {
        log.Fatal("启动服务器失败:", err)
    }
}
