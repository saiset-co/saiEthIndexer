package usecase

import (
	"context"

	"github.com/webmakom-com/saiBoilerplate/internal/types"
)

// Example struct
type SomeUseCase struct {
	repo SomeRepo
}

// New creates new usecase
func New(r SomeRepo) *SomeUseCase {
	return &SomeUseCase{
		repo: r,
	}
}

func (uc *SomeUseCase) GetAll(ctx context.Context) ([]*types.Some, error) {
	somes, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return somes, nil
}

func (uc *SomeUseCase) Set(ctx context.Context, some *types.Some) error {
	return uc.repo.Set(ctx, some)
}
