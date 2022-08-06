package tasks

import (
	"context"

	"github.com/webmakom-com/saiBoilerplate/types"
)

// Example struct
type Task struct {
	repo SomeRepo
}

// New creates new usecase
func New(r SomeRepo) *Task {
	return &Task{
		repo: r,
	}
}

func (t *Task) GetAll(ctx context.Context) ([]*types.Some, error) {
	somes, err := t.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return somes, nil
}

func (t *Task) Set(ctx context.Context, some *types.Some) error {
	return t.repo.Set(ctx, some)
}
