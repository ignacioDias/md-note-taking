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

func (r *Router) SetupRoutes() *http.ServeMux {

	r.mux.HandleFunc("POST /api/auth/register", r.sessionHandler.LoginUser)
	r.mux.HandleFunc("POST /api/auth/logout", r.sessionHandler.LogoutUser)
	r.mux.HandleFunc("POST /api/auth/login", r.sessionHandler.RegisterUser)

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
