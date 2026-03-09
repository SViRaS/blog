package models

import (
	"time"
)

type Post struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:255;not null"`
	Content   string `gorm:"type:text;not null"`
	Image     string `gorm:"size:255"`
	UserID    uint   `gorm:"not null;index"`
	User      User   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
