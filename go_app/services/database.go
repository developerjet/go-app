package services

import (
    "fmt"
    "go_app/config"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// DB 全局数据库实例
var DB *gorm.DB

// ConnectDB 初始化数据库连接
func ConnectDB() error {
    cfg, err := config.LoadConfig()
    if err != nil {
        return err
    }

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.DBName,
    )
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }

    DB = db
    return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
    return DB
}