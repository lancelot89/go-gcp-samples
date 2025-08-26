package service

import (
	"context"
	"testing"
	"time"

	"go-gcp-samples/v2-firestore/internal/firestore"
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
	if t.ID == "" {
		return firestore.ErrNotFound
	}
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

func TestTodoService_CreateTodo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test requiring Firestore client")
	}
	tests := []struct {
		name    string
		userID  string
		title   string
		wantErr bool
	}{
		{
			name:    "正常なTodo作成",
			userID:  "user123",
			title:   "テストタスク",
			wantErr: false,
		},
		{
			name:    "空のタイトル",
			userID:  "user123",
			title:   "",
			wantErr: false,
		},
		{
			name:    "空のユーザーID",
			userID:  "",
			title:   "テストタスク",
			wantErr: false,
		},
		{
			name:    "長いタイトル",
			userID:  "user123",
			title:   "これは非常に長いタイトルです。" + string(make([]byte, 1000)),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &firestore.TodoRepository{}
			svc := NewTodoService(repo)

			ctx := context.Background()
			todo, err := svc.CreateTodo(ctx, tt.userID, tt.title)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if todo == nil {
					t.Error("CreateTodo() returned nil todo without error")
					return
				}

				if todo.ID == "" {
					t.Error("CreateTodo() returned todo with empty ID")
				}

				if todo.UserID != tt.userID {
					t.Errorf("CreateTodo() userID = %v, want %v", todo.UserID, tt.userID)
				}

				if todo.Title != tt.title {
					t.Errorf("CreateTodo() title = %v, want %v", todo.Title, tt.title)
				}

				if todo.Done {
					t.Error("CreateTodo() done should be false by default")
				}

				if todo.CreatedAt.Before(time.Now().Add(-1*time.Minute)) || todo.CreatedAt.After(time.Now().Add(1*time.Minute)) {
					t.Error("CreateTodo() createdAt is not within expected range")
				}
			}
		})
	}
}

func TestTodoService_GetTodo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test requiring Firestore client")
	}
	ctx := context.Background()
	repo := &firestore.TodoRepository{}
	svc := NewTodoService(repo)

	tests := []struct {
		name    string
		id      string
		want    *firestore.Todo
		wantErr bool
	}{
		{
			name:    "存在しないTodo",
			id:      "nonexistent",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "空のID",
			id:      "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.GetTodo(ctx, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != nil {
				if got.ID != tt.want.ID {
					t.Errorf("GetTodo() ID = %v, want %v", got.ID, tt.want.ID)
				}
			}
		})
	}
}

func BenchmarkTodoService_CreateTodo(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark requiring Firestore client")
	}
	repo := &firestore.TodoRepository{}
	svc := NewTodoService(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.CreateTodo(ctx, "bench-user", "ベンチマークタスク")
	}
}