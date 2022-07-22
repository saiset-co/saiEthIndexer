package usecase

import (
	"context"

	"github.com/webmakom-com/saiBoilerplate/internal/entity"
)

type (
	SomeRepo interface {
		Set(ctx context.Context, entity *entity.Some) error
		GetAll(ctx context.Context) ([]*entity.Some, error)
	}
)
