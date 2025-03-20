package services

import (
    "fmt"
    "go_app/config"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
    cfg, err := config.LoadConfig()
    if err != nil {
        return nil, err
    }

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.DBName,
    )
    
    return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}