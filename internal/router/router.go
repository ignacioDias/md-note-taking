package router

import (
	"mdeditor/internal/handler"
	"mdeditor/internal/middleware"
	"net/http"
)

type Router struct {
	mux            *http.ServeMux
	userHandler    *handler.UserHandler
	sessionHandler *handler.SessionHandler
	noteHandler    *handler.NoteHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(userHandler *handler.UserHandler, sessionHandler *handler.SessionHandler, noteHandler *handler.NoteHandler, authMiddw *middleware.AuthMiddleware) *Router {
	return &Router{
		mux:            http.NewServeMux(),
		userHandler:    userHandler,
		sessionHandler: sessionHandler,
		noteHandler:    noteHandler,
		authMiddleware: authMiddw,
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {

	r.mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})
	r.mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/login.html")
	})
	r.mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/register.html")
	})
	r.mux.HandleFunc("GET /dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/dashboard.html")
	})
	r.mux.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/profile.html")
	})
	r.mux.HandleFunc("GET /settings", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/settings.html")
	})

	r.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	r.mux.HandleFunc("POST /api/auth/register", r.sessionHandler.RegisterUser)
	r.mux.HandleFunc("POST /api/auth/login", r.sessionHandler.LoginUser)
	r.mux.HandleFunc("DELETE /api/auth/logout", r.authMiddleware.AuthenticationMiddleware(r.sessionHandler.LogoutUser))

	r.mux.HandleFunc("GET /api/me", r.authMiddleware.AuthenticationMiddleware(r.userHandler.GetUser))
	r.mux.HandleFunc("PUT /api/me", r.authMiddleware.AuthenticationMiddleware(r.userHandler.UpdateUser))
	r.mux.HandleFunc("DELETE /api/me", r.authMiddleware.AuthenticationMiddleware(r.userHandler.DeleteUser))

	r.mux.HandleFunc("POST /api/notes", r.authMiddleware.AuthenticationMiddleware(r.noteHandler.CreateNote))
	r.mux.HandleFunc("GET /api/notes/{id}", r.authMiddleware.AuthenticationMiddleware(r.noteHandler.GetNote))
	r.mux.HandleFunc("DELETE /api/notes/{id}", r.authMiddleware.AuthenticationMiddleware(r.noteHandler.DeleteNote))
	r.mux.HandleFunc("PUT /api/notes/{id}", r.authMiddleware.AuthenticationMiddleware(r.noteHandler.UpdateNote))

	r.mux.HandleFunc("GET /api/me/notes", r.authMiddleware.AuthenticationMiddleware(r.noteHandler.GetNotesPerUser))

	return r.mux

}
