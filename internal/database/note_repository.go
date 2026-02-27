package database

import (
	"context"
	"database/sql"
	"errors"
	"mdeditor/internal/domain"

	"github.com/jmoiron/sqlx"
)

type NoteRepository struct {
	db *sqlx.DB
}

var ErrNoteNotFound error = errors.New("Note not found")

func (nRepo *NoteRepository) CreateNote(ctx context.Context, note *domain.Note) error {
	query := `
	INSERT INTO notes (user_id, title, content)
	VALUES (:user_id, :title, :content)
	RETURNING id, created_at, updated_at
	`
	stmt, err := nRepo.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, note, note)
}

func (nRepo *NoteRepository) FindNoteByIDAndUserID(ctx context.Context, userID, id int64) (*domain.Note, error) {
	var note domain.Note
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE id = $1 AND user_id = $2"
	if err := nRepo.db.GetContext(ctx, &note, query, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}
	return &note, nil
}

func (nRepo *NoteRepository) FindNotesByUserID(ctx context.Context, userID int64, limit, offset int) ([]domain.Note, error) {
	var notes []domain.Note
	query := "SELECT id, user_id, title, content, created_at, updated_at FROM notes WHERE user_id = $1 ORDER BY updated_at DESC LIMIT $2 OFFSET $3"
	if err := nRepo.db.SelectContext(ctx, &notes, query, userID, limit, offset); err != nil {
		return nil, err
	}
	return notes, nil
}

func (nRepo *NoteRepository) UpdateNoteByIDAndUserID(ctx context.Context, newNote *domain.Note, userID, id int64) error {
	query := `UPDATE notes SET title = $1, content = $2 WHERE id = $3 AND user_id = $4`
	result, err := nRepo.db.ExecContext(ctx, query, newNote.Title, newNote.Content, id, userID)
	return checkQueryResult(result, err)
}

func (nRepo *NoteRepository) DeleteNoteByIDAndUserID(ctx context.Context, userID, id int64) error {
	query := `DELETE FROM notes WHERE id = $1 AND user_id = $2`
	result, err := nRepo.db.ExecContext(ctx, query, id, userID)
	return checkQueryResult(result, err)
}

func checkQueryResult(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoteNotFound
	}
	return nil
}
