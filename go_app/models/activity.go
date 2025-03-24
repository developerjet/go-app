package models

import "time"

type Activity struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    StartTime   time.Time `json:"startTime"`
    EndTime     time.Time `json:"endTime"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
}