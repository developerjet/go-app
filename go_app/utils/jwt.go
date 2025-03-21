package utils

import (
    "time"
    "go_app/pkg/errcode"
    "github.com/golang-jwt/jwt"
)

type Claims struct {
    UserID       uint `json:"user_id"`
    TokenVersion int  `json:"token_version"`
    jwt.StandardClaims
}

func GenerateToken(userID uint, tokenVersion int) (string, error) {
    claims := Claims{
        UserID:       userID,
        TokenVersion: tokenVersion,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("your-secret-key")) // 请使用配置文件中的密钥
}

func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte("your-secret-key"), nil
    })

    if err != nil {
        ve, ok := err.(*jwt.ValidationError)
        if ok && ve.Errors&jwt.ValidationErrorExpired != 0 {
            return nil, errcode.TokenExpired
        }
        return nil, errcode.TokenInvalid
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errcode.TokenInvalid
    }

    return claims, nil
}