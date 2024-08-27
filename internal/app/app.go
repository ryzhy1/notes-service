package app

import (
	"github.com/go-chi/chi/v5"
	"log/slog"
	server "testovoe/internal/app/http"
	"testovoe/internal/handlers/notesHandlers"
	"testovoe/internal/routes"
	"testovoe/internal/services/notesService"
	"testovoe/internal/storage/postgres"
	"time"
)

type App struct {
	HTTPServer *server.Server
}

func New(log *slog.Logger, serverPort, storagePath string, tokenTTL time.Duration) *App {
	storage, err := postgres.New(storagePath)
	if err != nil {
		panic(err)
	}

	noteService := notesService.NewNotesService(log, storage)

	noteHandlers := notesHandlers.NewNotesHandlers(noteService)

	r := chi.NewRouter()
	r = routes.InitRoutes(log, noteHandlers, r)

	newServer := server.NewServer(log, serverPort, r)

	return &App{
		HTTPServer: newServer,
	}
}
