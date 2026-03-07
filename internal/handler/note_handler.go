package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"mdeditor/internal/database"
	"mdeditor/internal/domain"
	"mdeditor/internal/middleware"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type NoteHandler struct {
	noteRepo    *database.NoteRepository
	redisClient *redis.Client
}

type UpdateNote struct {
	Title   *string `json:"title" db:"title"`
	Content *string `json:"content" db:"content"`
}

func NewNoteHandler(noteRepo *database.NoteRepository, redisClient *redis.Client) *NoteHandler {
	return &NoteHandler{
		noteRepo:    noteRepo,
		redisClient: redisClient,
	}
}

func (noteHandler *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note domain.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	note.UserID = userID
	if err := noteHandler.noteRepo.CreateNote(r.Context(), &note); err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if json.NewEncoder(w).Encode(note) != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
	iter := noteHandler.redisClient.Scan(r.Context(), 0, fmt.Sprintf("notes:user:%d:*", userID), 0).Iterator()
	for iter.Next(r.Context()) {
		noteHandler.redisClient.Del(r.Context(), iter.Val())
	}
}
func (noteHandler *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
	idValue, ok := ParseIDParam(w, r)
	if !ok {
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	cacheKey := fmt.Sprintf("note:user:%d:id:%d", userID, idValue)
	cached, err := noteHandler.redisClient.Get(r.Context(), cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}
	note, err := noteHandler.noteRepo.FindNoteByIDAndUserID(r.Context(), userID, idValue)
	if ok := isActionValid(w, err, "find"); !ok {
		return
	}
	Encode(w, note)
	data, _ := json.Marshal(note)
	noteHandler.redisClient.Set(r.Context(), cacheKey, data, 5*time.Minute)
}
func (noteHandler *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idValue, ok := ParseIDParam(w, r)
	if !ok {
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	err := noteHandler.noteRepo.DeleteNoteByIDAndUserID(r.Context(), userID, idValue)
	if ok := isActionValid(w, err, "delete"); !ok {
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusNoContent)
	noteHandler.redisClient.Del(r.Context(), fmt.Sprintf("note:user:%d:id:%d", userID, idValue))
	iter := noteHandler.redisClient.Scan(r.Context(), 0, fmt.Sprintf("notes:user:%d:*", userID), 0).Iterator()
	for iter.Next(r.Context()) {
		noteHandler.redisClient.Del(r.Context(), iter.Val())
	}
}

func (noteHandler *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	idValue, ok := ParseIDParam(w, r)
	if !ok {
		return
	}
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var newNote UpdateNote
	if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	note, err := noteHandler.noteRepo.FindNoteByIDAndUserID(r.Context(), userID, idValue)
	if ok := isActionValid(w, err, "find"); !ok {
		return
	}
	if newNote.Title != nil {
		note.Title = *newNote.Title
	}
	if newNote.Content != nil {
		note.Content = *newNote.Content
	}
	err = noteHandler.noteRepo.UpdateNoteByIDAndUserID(r.Context(), note, userID, idValue)
	if ok := isActionValid(w, err, "update"); !ok {
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusOK)
	Encode(w, note)
	noteHandler.redisClient.Del(r.Context(), fmt.Sprintf("note:user:%d:id:%d", userID, idValue))
	iter := noteHandler.redisClient.Scan(r.Context(), 0, fmt.Sprintf("notes:user:%d:*", userID), 0).Iterator()
	for iter.Next(r.Context()) {
		noteHandler.redisClient.Del(r.Context(), iter.Val())
	}
}

func (noteHandler *NoteHandler) GetNotesPerUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	limit, offset := GetPaginationValues(r)

	cacheKey := fmt.Sprintf("notes:user:%d:limit:%d:offset:%d", userID, limit, offset)
	cached, err := noteHandler.redisClient.Get(r.Context(), cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	total, err := noteHandler.noteRepo.CountNotesByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to count notes", http.StatusInternalServerError)
		return
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	page := (offset / limit) + 1

	notes, err := noteHandler.noteRepo.FindNotesByUserID(r.Context(), userID, limit, offset)
	if ok := isActionValid(w, err, "find"); !ok {
		return
	}
	response := PaginatedResponse{
		Data:       notes,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
	Encode(w, response)
	data, _ := json.Marshal(response)
	noteHandler.redisClient.Set(r.Context(), cacheKey, data, 5*time.Minute)
}

func isActionValid(w http.ResponseWriter, err error, action string) bool {
	if err == nil {
		return true
	}

	if errors.Is(err, database.ErrNoteNotFound) {
		http.Error(w, "note not found", http.StatusNotFound)
		return false
	}

	http.Error(w, "failed to "+action+" note", http.StatusInternalServerError)
	return false
}
