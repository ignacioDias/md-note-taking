package handler

import (
	"html/template"
	"mdeditor/internal/database"
	"net/http"

	"golang.org/x/oauth2"
)

type UserHandler struct {
	userRepo *database.UserRepository
	config   *oauth2.Config
}

type userInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (userHandler *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
