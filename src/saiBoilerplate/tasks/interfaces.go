package tasks

import (
	"context"

	"github.com/webmakom-com/saiBoilerplate/types"
)

type (
	SomeRepo interface {
		Set(ctx context.Context, entity *types.Some) error
		GetAll(ctx context.Context) ([]*types.Some, error)
	}
)
