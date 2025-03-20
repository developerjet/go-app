package models

import "time"

type UserToken struct {
    ID        uint      `gorm:"primarykey"`
    UserID    uint      `gorm:"not null"`
    Token     string    `gorm:"type:varchar(255);not null"`
    Platform  string    `gorm:"type:varchar(20);not null"`
    ExpiredAt time.Time `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (ut *UserToken) IsExpired() bool {
    return time.Now().After(ut.ExpiredAt)
}