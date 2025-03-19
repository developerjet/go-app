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
        if token == "" {
            c.JSON(http.StatusUnauthorized, models.Response{Error: "未提供认证token"})
            c.Abort()
            return
        }

        claims, err := utils.ParseToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, models.Response{Error: "无效的token"})
            c.Abort()
            return
        }

        c.Set("userID", claims.UserID)
        c.Next()
    }
}
