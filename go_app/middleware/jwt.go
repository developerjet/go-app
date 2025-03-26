package middleware

import (
	"go_app/models"
	"go_app/pkg/errcode"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 打印请求信息
        log.Printf("请求路径: %s", c.Request.URL.Path)
        log.Printf("请求方法: %s", c.Request.Method)
        
        // 打印所有请求头
        for k, v := range c.Request.Header {
            log.Printf("Header: %s = %v", k, v)
        }

        token := c.Query("token")
        if token == "" {
            token = c.GetHeader("Authorization")
            log.Printf("从 Header 获取到的 token: %s", token)
        }
        
        if token == "" {
            log.Printf("未找到 token，请检查 Authorization header 是否正确设置")
            c.JSON(http.StatusOK, models.NewError(errcode.Unauthorized))
            c.Abort()
            return
        }

        userID, err := models.ParseToken(token)
        if err != nil {
            log.Printf("token 解析失败，错误: %v\ntoken 值: %s", err, token)
            c.JSON(http.StatusOK, models.NewError(errcode.Unauthorized))
            c.Abort()
            return
        }

        log.Printf("token 验证成功，用户 ID: %d", userID)
        c.Set("userId", userID)
        c.Next()
    }
}