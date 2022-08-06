package socket

import (
	"context"

	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/tasks"
	"go.uber.org/zap"
)

type Handler struct {
	notify chan error
	logger *zap.Logger
	cfg    *config.Configuration
	task   *tasks.Task
}

func New(ctx context.Context, cfg *config.Configuration, logger *zap.Logger, t *tasks.Task) *Handler {
	s := &Handler{
		notify: make(chan error, 1),
		cfg:    cfg,
		logger: logger,
		task:   t,
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
