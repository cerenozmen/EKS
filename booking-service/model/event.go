package model

import "time"

type Event struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"column:name"`
	UserId    int       `json:"user_id" gorm:"column:user_id"`
	Capacity  int       `json:"capacity"`
	IsActive  bool      `json:"isActive" gorm:"column:is_active"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}
