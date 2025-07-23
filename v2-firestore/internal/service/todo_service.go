package service

import (
	"context"
	"fmt"

	"go-gcp-samples/v2-firestore/internal/firestore"
)

type TodoService struct {
	repo *firestore.TodoRepository
}

func NewTodoService(repo *firestore.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(ctx context.Context, userID, title string) (*firestore.Todo, error) {
	// In a real application, you might generate a unique ID here
	todo := &firestore.Todo{
		ID:     fmt.Sprintf("%s-%s", userID, title), // Simple ID for demonstration
		UserID: userID,
		Title:  title,
	}

	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}
	return todo, nil
}

func (s *TodoService) GetTodo(ctx context.Context, id string) (*firestore.Todo, error) {
	todo, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	return todo, nil
}
