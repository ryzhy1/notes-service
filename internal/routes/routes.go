package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/oauth"
	"log/slog"
	"testovoe/internal/handlers/notesHandlers"
	oa "testovoe/internal/lib/oauth"
	"testovoe/internal/middlewares"
)

func InitRoutes(log *slog.Logger, notesHandlers *notesHandlers.NotesHandlers, router *chi.Mux) *chi.Mux {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middlewares.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "HEAD", "OPTION"},
		AllowedHeaders:   []string{"User-Agent", "Content-Type", "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Origin", "Pragma", "Referer"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	oa.AuthAPI(router)
	registerAPI(router, notesHandlers)

	return router
}

func registerAPI(r *chi.Mux, notesHandlers *notesHandlers.NotesHandlers) {
	r.Route("/", func(r chi.Router) {
		// use the Bearer Authentication middleware
		r.Use(oauth.Authorize("yaroslav-the-best", nil))
		r.Post("/add-note", notesHandlers.AddNote)
		r.Get("/get-notes", notesHandlers.GetNotes)
	})
}
