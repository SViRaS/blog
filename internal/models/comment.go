package models

import "time"

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"type:text;not null"`
	UserID    uint   `gorm:"not null;index"`
	User      User   `gorm:"foreignKey:UserID"`
	PostID    uint   `gorm:"not null;index"`
	Post      Post   `gorm:"foreignKey:PostID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
