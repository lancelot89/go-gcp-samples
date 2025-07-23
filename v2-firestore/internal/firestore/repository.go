package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

type TodoRepository struct {
	cli *firestore.Client
}

func NewTodoRepository(cli *firestore.Client) *TodoRepository {
	return &TodoRepository{cli: cli}
}

func (r *TodoRepository) Create(ctx context.Context, t *Todo) error {
	_, err := r.cli.Collection("todos").Doc(t.ID).Set(ctx, t)
	return err
}

func (r *TodoRepository) Get(ctx context.Context, id string) (*Todo, error) {
	snap, err := r.cli.Collection("todos").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var t Todo
	if err := snap.DataTo(&t); err != nil {
		return nil, err
	}
	return &t, nil
}
