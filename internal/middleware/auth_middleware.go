package middleware

import (
	"context"
	"mdeditor/internal/database"
	"net/http"
)

type contextKey string

type AuthMiddleware struct {
	sessionRepo *database.SessionRepository
}

const userIDKey contextKey = "userID"

func NewAuthMiddleware(sessionRepo *database.SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo: sessionRepo,
	}
}

func (auth *AuthMiddleware) AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		session, err := auth.sessionRepo.FindSessionByID(r.Context(), cookie.Value)
		if err == database.ErrSessionNotFound {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, session.UserID)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func GetUserID(r *http.Request) (int64, bool) {
	userID, ok := r.Context().Value(userIDKey).(int64)
	return userID, ok
}
