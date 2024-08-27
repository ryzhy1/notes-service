package notesHandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testovoe/internal/models"
	"testovoe/internal/services/notesService"
)

type NotesService interface {
	AddNote(ctx context.Context, content, owner string) (string, error)
	GetNotes(ctx context.Context, owner string) ([]models.Note, error)
}

type NotesHandlers struct {
	service NotesService
}

func NewNotesHandlers(service *notesService.NotesService) *NotesHandlers {
	return &NotesHandlers{
		service: service,
	}
}

func (h *NotesHandlers) AddNote(w http.ResponseWriter, r *http.Request) {
	var noteReq models.Note

	err := json.NewDecoder(r.Body).Decode(&noteReq)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	username, err := ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	noteID, err := h.service.AddNote(r.Context(), noteReq.Content, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	noteResp := models.Note{
		ID:      noteID,
		Content: noteReq.Content,
		Owner:   username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(noteResp)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *NotesHandlers) GetNotes(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	username, err := ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	notes, err := h.service.GetNotes(r.Context(), username)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		return
	}
}

func ValidateToken(token string) (string, error) {
	return "user1", nil
}
