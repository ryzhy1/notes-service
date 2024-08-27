package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"testovoe/internal/models"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(conn string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := pgxpool.New(context.Background(), conn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) AddNote(noteId, content, owner string) (string, error) {
	_, err := s.db.Exec(context.Background(),
		`INSERT INTO notes (id, content, owner) 
			VALUES ($1, $2, $3)`,
		noteId, content, owner)

	return noteId, err
}

func (s *Storage) GetNotes(owner string) ([]models.Note, error) {
	notes := make([]models.Note, 0)

	rows, err := s.db.Query(context.Background(),
		`SELECT id, content, owner
			FROM notes
			WHERE owner = $1`,
		owner)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		note := models.Note{}
		err = rows.Scan(&note.ID, &note.Content, &note.Owner)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *Storage) Close() {
	s.db.Close()

	return
}
