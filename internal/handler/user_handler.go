package handler

import (
	"encoding/json"
	"errors"
	"mdeditor/internal/database"
	"mdeditor/internal/middleware"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo *database.UserRepository
}

type UserUpdateRequest struct {
	Email          *string `json:"email"`
	Name           *string `json:"name"`
	OldPassword    string  `json:"oldPassword"`
	NewPassword    *string `json:"newPassword"`
	ProfilePicture *string `json:"profilePicture"`
}

func NewUserHandler(userRepo *database.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (userHandler *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := userHandler.userRepo.FindUserByID(r.Context(), userID)
	if ok := isUserFound(w, err); ok {
		Encode(w, &user)
	}
}

func (userHandler *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := userHandler.userRepo.FindUserByID(r.Context(), userID)
	if ok := isUserFound(w, err); !ok {
		return
	}

	var userUpdateReq UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&userUpdateReq); err != nil {
		http.Error(w, "Invalid input to update", http.StatusBadRequest)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userUpdateReq.OldPassword)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if userUpdateReq.Email != nil && *userUpdateReq.Email != "" {
		if IsValidEmail(*userUpdateReq.Email) {
			user.Email = *userUpdateReq.Email
		} else {
			http.Error(w, "Invalid email to update", http.StatusBadRequest)
			return
		}
		if isUsed, err := userHandler.userRepo.IsMailUsed(r.Context(), *userUpdateReq.Email); isUsed || err != nil {
			http.Error(w, "Email repeated", http.StatusBadRequest)
			return
		}
	}

	if userUpdateReq.NewPassword != nil && *userUpdateReq.NewPassword != "" {
		if IsValidPassword(*userUpdateReq.NewPassword) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*userUpdateReq.NewPassword), 14)
			if err != nil {
				http.Error(w, "Failed to process password", http.StatusInternalServerError)
				return
			}
			user.HashedPassword = string(hashedPassword)
		} else {
			http.Error(w, "Invalid password to update", http.StatusBadRequest)
			return
		}
	}

	if userUpdateReq.Name != nil && *userUpdateReq.Name != "" {
		user.Name = *userUpdateReq.Name
	}
	if userUpdateReq.ProfilePicture != nil && *userUpdateReq.ProfilePicture != "" {
		user.ProfilePicture = *userUpdateReq.ProfilePicture
	}

	if err := userHandler.userRepo.UpdateUserByID(r.Context(), userID, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Encode(w, &user)

}

func (userHandler *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := userHandler.userRepo.DeleteUserByID(r.Context(), userID); err != nil {
		if err == database.ErrUserNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func isUserFound(w http.ResponseWriter, err error) bool {
	if errors.Is(err, database.ErrUserNotFound) {
		http.Error(w, "User not found", http.StatusNotFound)
		return false
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	return true
}
