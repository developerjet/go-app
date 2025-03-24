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
        // 支持从 URL 参数获取 token
        token := c.Query("token")
        if token == "" {
            token = c.GetHeader("Authorization")
        }
        
        if token == "" {
            log.Printf("未找到 token")
            c.JSON(http.StatusOK, models.NewError(errcode.Unauthorized))
            c.Abort()
            return
        }

        // 直接使用 token，不再检查和移除 Bearer 前缀
        userID, err := models.ParseToken(token)
        if err != nil {
            log.Printf("token 解析失败: %v", err)
            c.JSON(http.StatusOK, models.NewError(errcode.Unauthorized))
            c.Abort()
            return
        }

        c.Set("userId", userID)
        c.Next()
    }
}