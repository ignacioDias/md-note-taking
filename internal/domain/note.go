package domain

import "time"

type Note struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"userId" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func NewNote(userID int64, title string, content string) *Note {
	return &Note{
		UserID:  userID,
		Title:   title,
		Content: content,
	}
}
