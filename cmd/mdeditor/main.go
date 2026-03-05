package main

import (
	"mdeditor/internal/database"
	"mdeditor/internal/handler"
	"mdeditor/internal/middleware"
	"mdeditor/internal/router"
	"mdeditor/internal/server"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"
)

func main() {
	godotenv.Load()
	redisClient := database.NewRedisClient(os.Getenv("REDIS_URL"))
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	database := database.NewDatabase(db)
	if err := database.Init(); err != nil {
		panic(err)
	}

	userHandler := handler.NewUserHandler(database.UserRepo)
	sessionHandler := handler.NewSessionHandler(database.UserRepo, database.SessionRepo)
	noteHandler := handler.NewNoteHandler(database.NoteRepo, redisClient)
	authMid := middleware.NewAuthMiddleware(database.SessionRepo)
	router := router.NewRouter(userHandler, sessionHandler, noteHandler, authMid)

	server := server.NewServer(os.Getenv("PORT"), router)
	server.StartServer()
}
