package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gcp-samples/v2-firestore/internal/firestore"
	"go-gcp-samples/v2-firestore/internal/service"
)

type mockTodoRepository struct {
	todos map[string]*firestore.Todo
}

func newMockTodoRepository() *mockTodoRepository {
	return &mockTodoRepository{
		todos: make(map[string]*firestore.Todo),
	}
}

func (m *mockTodoRepository) Create(ctx context.Context, t *firestore.Todo) error {
	m.todos[t.ID] = t
	return nil
}

func (m *mockTodoRepository) Get(ctx context.Context, id string) (*firestore.Todo, error) {
	todo, ok := m.todos[id]
	if !ok {
		return nil, firestore.ErrNotFound
	}
	return todo, nil
}

func TestTodoHandler_CreateTodo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test requiring Firestore client")
	}

	tests := []struct {
		name       string
		request    CreateTodoRequest
		wantStatus int
		wantErr    bool
	}{
		{
			name: "正常なTodo作成",
			request: CreateTodoRequest{
				UserID: "user123",
				Title:  "テストTodo",
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "タイトルが空",
			request: CreateTodoRequest{
				UserID: "user123",
				Title:  "",
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "ユーザーIDが空",
			request: CreateTodoRequest{
				UserID: "",
				Title:  "テストTodo",
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &firestore.TodoRepository{}
			svc := service.NewTodoService(repo)
			handler := NewTodoHandler(svc)

			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.CreateTodo(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("CreateTodo() status = %v, want %v", rec.Code, tt.wantStatus)
			}

			if tt.wantErr {
				return
			}

			var result firestore.Todo
			if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
				t.Errorf("CreateTodo() failed to decode response: %v", err)
			}

			if result.Title != tt.request.Title {
				t.Errorf("CreateTodo() title = %v, want %v", result.Title, tt.request.Title)
			}
		})
	}
}

func TestTodoHandler_GetTodo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test requiring Firestore client")
	}
	tests := []struct {
		name       string
		queryID    string
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "IDパラメータなし",
			queryID:    "",
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "存在しないID",
			queryID:    "nonexistent",
			wantStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &firestore.TodoRepository{}
			svc := service.NewTodoService(repo)
			handler := NewTodoHandler(svc)

			req := httptest.NewRequest(http.MethodGet, "/todos?id="+tt.queryID, nil)
			rec := httptest.NewRecorder()

			handler.GetTodo(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("GetTodo() status = %v, want %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestTodoHandler_CreateTodo_InvalidJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test requiring Firestore client")
	}
	repo := &firestore.TodoRepository{}
	svc := service.NewTodoService(repo)
	handler := NewTodoHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateTodo(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("CreateTodo() with invalid JSON: status = %v, want %v", rec.Code, http.StatusBadRequest)
	}
}
