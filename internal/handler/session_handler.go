package handler

import (
	"encoding/json"
	"html/template"
	"mdeditor/internal/database"
	"mdeditor/internal/domain"
	"mdeditor/internal/middleware"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type SessionHandler struct {
	userRepo    *database.UserRepository
	sessionRepo *database.SessionRepository
}

type UserRegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

var DEFAULT_PIC string = "https://origin.giantbomb.com/a/uploads/scale_medium/8/82962/1543003-adventure_time_with_finn_and_jake_john_dimaggio_2.jpg"

func (sessionHandler *SessionHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userRequest UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}
	if userRequest.Name == "" {
		http.Error(w, "Invalid name: can't be empty", http.StatusBadRequest)
		return
	}
	if !IsValidEmail(userRequest.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	if !IsValidPassword(userRequest.Password) {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 14)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}
	user := domain.NewUser(userRequest.Email, userRequest.Name, string(hashedPassword), DEFAULT_PIC)
	if err := sessionHandler.userRepo.CreateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (sessionHandler *SessionHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
func (sessionHandler *SessionHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if sessionHandler.sessionRepo.DeleteSessionsByUserID(r.Context(), userID) != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
