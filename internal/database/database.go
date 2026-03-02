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
		NoteRepo: &NoteRepository{db: db},
	}
}

var createNotesTable = `
CREATE TABLE IF NOT EXISTS notes (
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_notes_user_created
ON notes (user_id, created_at DESC);
`
var createUsersTable string = `
CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    email       TEXT UNIQUE NOT NULL,
    name        TEXT,
    created_at  TIMESTAMP DEFAULT NOW(),
    hashed_password TEXT NOT NULL,
    profile_picture TEXT
);`

var createSessionsTable string = `
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,              
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);`

func (d *Database) Init() error {
	_, err := d.DB.Exec(createUsersTable)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec(createNotesTable)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec(createSessionsTable)
	if err != nil {
		return err
	}
	return nil
}
