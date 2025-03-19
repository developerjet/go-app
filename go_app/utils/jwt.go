package utils

import (
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v4"  // 更新为 v4 版本
)

// Claims 自定义的 JWT Claims
type Claims struct {
    UserID uint `json:"userId"`  // 添加 json 标签
    jwt.StandardClaims
}

var jwtSecret = []byte("your_jwt_secret_key")

func GenerateToken(userID uint) (string, error) {
    claims := &Claims{
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret) // 使用 jwtSecret 替换 secretKey
}

func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil // 使用 jwtSecret 替换 secretKey
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}