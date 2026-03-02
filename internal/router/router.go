package router

import (
	"mdeditor/internal/handler"
	"net/http"
)

type Router struct {
	mux         *http.ServeMux
	userHandler *handler.UserHandler
	noteHandler *handler.NoteHandler
}

func (r *Router) SetupRoutes() *http.ServeMux {

	r.mux.HandleFunc("GET /auth/login", r.userHandler.LoginHandler)

	return r.mux

}
