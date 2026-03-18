package models

import "time"

type Payload struct {
	Event          string    `json:"event"`
	PostID         uint      `json:"post_id"`
	CommentID      uint      `json:"comment_id"`
	AuthorID       uint      `json:"author_id"`
	AuthorUsername string    `json:"author_username"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}
