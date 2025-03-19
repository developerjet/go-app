package utils

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("your_jwt_secret_key")

func GenerateToken(userID uint) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// 添加验证 token 的函数
func ParseToken(tokenString string) (uint, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        return 0, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID := uint(claims["user_id"].(float64))
        return userID, nil
    }

    return 0, errors.New("invalid token")
}