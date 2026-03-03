package handler

import (
	"encoding/json"
	"html/template"
	"mdeditor/internal/database"
	"mdeditor/internal/domain"
	"mdeditor/internal/middleware"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo    *database.UserRepository
	sessionRepo *database.SessionRepository
}

type UserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

var DEFAULT_PIC string = "https://origin.giantbomb.com/a/uploads/scale_medium/8/82962/1543003-adventure_time_with_finn_and_jake_john_dimaggio_2.jpg"

func (userHandler *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userRequest UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}
	if userRequest.Name == "" {
		http.Error(w, "Invalid name: can't be empty", http.StatusBadRequest)
		return
	}
	if !isValidEmail(userRequest.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	if !isValidPassword(userRequest.Password) {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), 14)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}
	user := domain.NewUser(userRequest.Email, userRequest.Name, string(hashedPassword), DEFAULT_PIC)
	if err := userHandler.userRepo.CreateUser(r.Context(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (userHandler *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}
func (userHandler *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if userHandler.sessionRepo.DeleteSessionsByUserID(r.Context(), userID) != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 72 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}
