package models

import "time"

type Notification struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"userId"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Read      bool      `json:"read" gorm:"default:false"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}