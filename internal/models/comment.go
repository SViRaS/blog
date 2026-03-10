package models

import "time"

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:255;not null"`
	Image     string `gorm:"size:255"`
	UserID    uint   `gorm:"not null;index"`
	User      User   `gorm:"foreignKey:UserID"`
	PostID    uint   `gorm:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
