package database

import (
	"database/sql"
	"errors"
	"mdeditor/internal/domain"

	"github.com/jmoiron/sqlx"
)

type NoteRepository struct {
	db *sqlx.DB
}

var ErrNoteNotFound error = errors.New("Note not found")

func (nRepo *NoteRepository) CreateNote(note *domain.Note) error {
	query := `
	INSERT INTO notes (user_id, title, content)
	VALUES (:user_id, :title, :content)
	RETURNING id, created_at, updated_at
	`
	return nRepo.db.Get(note, query, note)
}

func (nRepo *NoteRepository) FindNoteByIDAndUserID(userID, id int64) (*domain.Note, error) {
	var note domain.Note
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2"
	if err := nRepo.db.Get(&note, query, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}
	return &note, nil
}

func (nRepo *NoteRepository) FindNotesByUserID(userID int64, limit, offset int) ([]domain.Note, error) {
	var notes []domain.Note
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE user_id = $1 ORDER BY updated_at DESC LIMIT $2 OFFSET $3"
	if err := nRepo.db.Select(&notes, query, userID, limit, offset); err != nil {
		return nil, err
	}
	return notes, nil
}
