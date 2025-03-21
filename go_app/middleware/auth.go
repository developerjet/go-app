package middleware

import (
	"go_app/models"
	"go_app/services"
    "go_app/pkg/errcode"
	"go_app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusOK, models.NewError(errcode.TokenMissing))
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, models.NewError(errcode.TokenInvalid))
			c.Abort()
			return
		}

		userID := claims.UserID
		tokenVersion := claims.TokenVersion

		// 验证 token 版本
		db := services.GetDB()
		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusOK, models.NewError(errcode.UserNotFound))
			c.Abort()
			return
		}

		if tokenVersion != user.TokenVersion {
			c.JSON(http.StatusOK, models.NewError(errcode.TokenVersionError))
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
