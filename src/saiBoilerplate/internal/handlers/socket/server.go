package socket

import (
	"context"

	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase"
	"go.uber.org/zap"
)

type Handler struct {
	notify  chan error
	logger  *zap.Logger
	cfg     *config.Configuration
	uc      *usecase.SomeUseCase
	handler *Handler
}

func New(ctx context.Context, cfg *config.Configuration, logger *zap.Logger, uc *usecase.SomeUseCase) *Handler {
	s := &Handler{
		notify: make(chan error, 1),
		cfg:    cfg,
		logger: logger,
		uc:     uc,
	}
	s.start(ctx)
	return s
}

func (s *Handler) start(ctx context.Context) {
	go func() {
		s.notify <- s.socketStart(ctx)
		close(s.notify)
	}()
}

func (s *Handler) Notify() <-chan error {
	return s.notify
}
