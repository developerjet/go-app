package services

import (
    "log"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "go_app/models"
)

func ConnectDB() (*gorm.DB, error) {
    dsn := "root:123456@tcp(127.0.0.1:3306)/go_app?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
        return nil, err
    }
    
    // 自动迁移
    db.AutoMigrate(&models.User{})
    
    return db, nil
}