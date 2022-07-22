package usecase

import (
	"context"

	"github.com/webmakom-com/saiBoilerplate/internal/entity"
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

func (uc *SomeUseCase) GetAll(ctx context.Context) ([]*entity.Some, error) {
	somes, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return somes, nil
}
