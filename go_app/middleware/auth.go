package middleware

import (
    "go_app/models"
    "go_app/utils"
    "net/http"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // 修改错误响应格式
        if token == "" {
            c.JSON(http.StatusOK, models.NewError("请提供认证令牌"))
            c.Abort()
            return
        }

        claims, err := utils.ParseToken(token)
        if err != nil {
            c.JSON(http.StatusOK, models.NewError("无效的认证令牌"))
            c.Abort()
            return
        }

        c.Set("userID", claims.UserID)
        c.Next()
    }
}
