package usecase

import (
	"context"

	"github.com/webmakom-com/saiBoilerplate/internal/types"
)

type (
	SomeRepo interface {
		Set(ctx context.Context, entity *types.Some) error
		GetAll(ctx context.Context) ([]*types.Some, error)
	}
)
