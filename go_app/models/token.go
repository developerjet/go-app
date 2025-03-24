package models

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key")

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (uint, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil  // 使用 jwtSecret 替代 SecretKey
    })
    
    if err != nil {
        return 0, err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := uint(claims["user_id"].(float64))
        return userID, nil
    }
    
    return 0, errors.New("invalid token claims")
}