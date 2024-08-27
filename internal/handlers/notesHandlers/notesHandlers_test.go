package notesHandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testovoe/internal/models"
)

// MockNotesService - простой мок для NotesService
type MockNotesService struct {
	addNoteFunc  func(ctx context.Context, content, owner string) (string, error)
	getNotesFunc func(ctx context.Context, owner string) ([]models.Note, error)
}

func (m *MockNotesService) AddNote(ctx context.Context, content, owner string) (string, error) {
	return m.addNoteFunc(ctx, content, owner)
}

func (m *MockNotesService) GetNotes(ctx context.Context, owner string) ([]models.Note, error) {
	return m.getNotesFunc(ctx, owner)
}

func TestNotesHandlers_AddNote(t *testing.T) {
	tests := []struct {
		name          string
		service       NotesService
		requestBody   string
		expectedCode  int
		expectedBody  models.Note
		authHeader    string
		validateToken func(string) (string, error)
	}{
		{
			name: "Valid Request",
			service: &MockNotesService{
				addNoteFunc: func(ctx context.Context, content, owner string) (string, error) {
					return "123e4567-e89b-12d3-a456-426614174000", nil // UUID format
				},
			},
			requestBody:  `{"content":"Test note"}`,
			expectedCode: http.StatusCreated,
			expectedBody: models.Note{Content: "Test note", Owner: "user1"},
			authHeader:   "Bearer valid-token",
			validateToken: func(token string) (string, error) {
				return "user1", nil
			},
		},
		{
			name:         "Invalid JSON",
			service:      &MockNotesService{},
			requestBody:  `{"content":`,
			expectedCode: http.StatusBadRequest,
			expectedBody: models.Note{},
			authHeader:   "Bearer valid-token",
			validateToken: func(token string) (string, error) {
				return "user1", nil
			},
		},
		{
			name:         "Missing Authorization",
			service:      &MockNotesService{},
			requestBody:  `{"content":"Test note"}`,
			expectedCode: http.StatusUnauthorized,
			expectedBody: models.Note{},
			authHeader:   "",
			validateToken: func(token string) (string, error) {
				return "", errors.New("invalid token")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/add-note", bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Set("Authorization", tt.authHeader)
			w := httptest.NewRecorder()

			h := &NotesHandlers{service: tt.service}
			h.AddNote(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("status code = %v, want %v", resp.StatusCode, tt.expectedCode)
			}

			if resp.StatusCode == http.StatusCreated {
				var responseBody models.Note
				if err := json.Unmarshal(body, &responseBody); err != nil {
					t.Errorf("Failed to unmarshal response body: %v", err)
					return
				}

				if uuid.Validate(responseBody.ID) != nil {
					t.Errorf("ID = %v does not match UUID format", responseBody.ID)
				}
				if responseBody.Content != tt.expectedBody.Content || responseBody.Owner != tt.expectedBody.Owner {
					t.Errorf("body = %+v, want %+v", responseBody, tt.expectedBody)
				}
			}
		})
	}
}

func TestNotesHandlers_GetNotes(t *testing.T) {
	tests := []struct {
		name          string
		service       NotesService
		expectedCode  int
		expectedBody  []models.Note
		authHeader    string
		validateToken func(string) (string, error)
	}{
		{
			name: "Valid Request",
			service: &MockNotesService{
				getNotesFunc: func(ctx context.Context, owner string) ([]models.Note, error) {
					return []models.Note{
						{ID: "123e4567-e89b-12d3-a456-426614174000", Content: "Test note", Owner: "user1"},
					}, nil
				},
			},
			expectedCode: http.StatusOK,
			expectedBody: []models.Note{
				{Content: "Test note", Owner: "user1"},
			},
			authHeader: "Bearer valid-token",
			validateToken: func(token string) (string, error) {
				return "user1", nil
			},
		},
		{
			name:         "Missing Authorization",
			service:      &MockNotesService{},
			expectedCode: http.StatusUnauthorized,
			expectedBody: []models.Note{}, // Пустое тело
			authHeader:   "",
			validateToken: func(token string) (string, error) {
				return "", errors.New("invalid token")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/get-notes", nil)
			req.Header.Set("Authorization", tt.authHeader)
			w := httptest.NewRecorder()

			h := &NotesHandlers{service: tt.service}
			h.GetNotes(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("status code = %v, want %v", resp.StatusCode, tt.expectedCode)
			}

			if resp.StatusCode == http.StatusOK {
				var responseBody []models.Note
				if err := json.Unmarshal(body, &responseBody); err != nil {
					t.Errorf("Failed to unmarshal response body: %v", err)
					return
				}

				for _, note := range responseBody {
					if uuid.Validate(note.ID) != nil {
						t.Errorf("ID = %v does not match UUID format", note.ID)
					}
				}
				if len(responseBody) != len(tt.expectedBody) {
					t.Errorf("body length = %v, want %v", len(responseBody), len(tt.expectedBody))
					return
				}
				for i, note := range responseBody {
					if note.Content != tt.expectedBody[i].Content || note.Owner != tt.expectedBody[i].Owner {
						t.Errorf("body[%d] = %+v, want %+v", i, note, tt.expectedBody[i])
					}
				}
			}
		})
	}
}

func Test_validateToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid token",
			token:   "valid-token",
			want:    "user1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
