package handler

import (
	"html/template"
	"mdeditor/internal/database"
	"net/http"
)

type UserHandler struct {
	userRepo    *database.UserRepository
	sessionRepo *database.SessionRepository
}

func (userHandler *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
