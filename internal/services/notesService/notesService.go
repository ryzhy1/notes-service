package notesService

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"testovoe/internal/middlewares"
	"testovoe/internal/models"
	spellcheck "testovoe/internal/services/spellchecker"
	"testovoe/internal/storage/postgres"
)

type NotesStorage interface {
	AddNote(noteId, content, owner string) (string, error)
	GetNotes(owner string) ([]models.Note, error)
	Close()
}

type NotesService struct {
	log *slog.Logger
	db  NotesStorage
}

func NewNotesService(log *slog.Logger, db *postgres.Storage) *NotesService {
	return &NotesService{
		log: log,
		db:  db,
	}
}

func (s *NotesService) AddNote(ctx context.Context, content, owner string) (string, error) {
	const op = "notesService.AddNote"

	s.log.With(
		slog.String("op", op),
		slog.String(owner, owner),
		slog.String("request id", middleware.GetReqID(ctx)),
	)

	s.log.Info("checking if content has spelling errors", slog.String("owner", owner))

	spellErrors, err := spellcheck.CheckSpelling(content)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if len(spellErrors) > 0 {
		return "", fmt.Errorf("content has spelling errors")
	}

	s.log.Info("creating uuid for note", slog.String("owner", owner))

	noteId, err := middlewares.UUIDGenerator()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("adding note", slog.String("owner", owner))

	noteID, err := s.db.AddNote(noteId.String(), content, owner)
	if err != nil {
		s.log.Error("failed to add note to the database", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("note added", slog.String("owner", owner))

	return noteID, nil
}

func (s *NotesService) GetNotes(ctx context.Context, owner string) ([]models.Note, error) {
	const op = "notesService.GetNotes"

	s.log.With(
		slog.String("op", op),
		slog.String(owner, owner),
		slog.String("request id", middleware.GetReqID(ctx)),
	)

	s.log.Info("getting notes", slog.String("owner", owner))

	notes, err := s.db.GetNotes(owner)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("got notes", slog.String("owner", owner))

	return notes, nil
}

func (s *NotesService) Close() {
	s.db.Close()
}
