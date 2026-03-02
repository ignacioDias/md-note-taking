package handler

import "mdeditor/internal/database"

type NoteHandler struct {
	noteRepo *database.NoteRepository
}
