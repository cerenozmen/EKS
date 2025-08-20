package model

import "time"

type Booking struct {
    Id         int       `json:"id" gorm:"primaryKey"`
    UserId     int       `json:"user_id" gorm:"user_id"`
	EventId    int       `json:"event_id" gorm:"event_id"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}
