package models

import (
	"time"
)

type Post struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:255;not null"`
	Content   string `gorm:"type:text;not null"`
	UserID    uint   `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
