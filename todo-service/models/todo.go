package models

import (
	"gorm.io/gorm"
	"time"
)

type Todo struct {
	gorm.Model
	UserID      uint      `gorm:"not null" json:"user_id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"size:1000" json:"description"`
	Completed   bool      `gorm:"default:false" json:"completed"`
	DueDate     time.Time `json:"due_date"`
}
