package database

import "github.com/jmoiron/sqlx"

type Database struct {
	DB       *sqlx.DB
	UserRepo *UserRepository
	NoteRepo *NoteRepository
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		DB:       db,
		UserRepo: &UserRepository{db: db},
	}
}

var createNotesTable = `
CREATE TABLE IF NOT EXISTS notes (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notes_user_created
ON notes (user_id, created_at DESC);
`
var createUsersTable string
