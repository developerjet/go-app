package middleware

import (
    "go_app/models"
    "go_app/services"
    "go_app/utils"
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusOK, models.NewTokenError(models.TokenStatusInvalid))
            c.Abort()
            return
        }

        // 先验证 JWT token 的基本格式
        claims, err := utils.ParseToken(token)
        if err != nil {
            c.JSON(http.StatusOK, models.NewTokenError(models.TokenStatusInvalid))
            c.Abort()
            return
        }

        db := services.GetDB()
        
        // 从 UserToken 表中检查 token 状态
        var userToken models.UserToken
        if err := db.Where("user_id = ? AND token = ?", claims.UserID, token).First(&userToken).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                // 检查是否存在其他有效token，说明是在其他设备登录
                var otherToken models.UserToken
                if err := db.Where("user_id = ?", claims.UserID).First(&otherToken).Error; err == nil {
                    c.JSON(http.StatusOK, models.NewTokenError(models.TokenStatusReplace))
                } else {
                    c.JSON(http.StatusOK, models.NewTokenError(models.TokenStatusInvalid))
                }
            } else {
                c.JSON(http.StatusOK, models.NewError("验证token失败"))
            }
            c.Abort()
            return
        }

        // 检查 token 是否过期
        if userToken.IsExpired() {
            // 删除过期的token
            db.Delete(&userToken)
            c.JSON(http.StatusOK, models.NewTokenError(models.TokenStatusExpired))
            c.Abort()
            return
        }

        // 设置用户信息到上下文
        c.Set("userID", claims.UserID)
        c.Next()
    }
}