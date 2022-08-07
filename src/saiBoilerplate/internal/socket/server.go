package socket

import (
	"context"
	"net"
	"time"

	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/tasks"
	"go.uber.org/zap"
)

const (
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	notify   chan error
	listener net.Listener
	logger   *zap.Logger
	cfg      *config.Configuration
	task     *tasks.Task
}

func New(ctx context.Context, cfg *config.Configuration, logger *zap.Logger, t *tasks.Task, errChan chan error) (*Server, error) {
	ln, err := net.Listen("tcp", cfg.Common.SocketServer.Host+":"+cfg.Common.SocketServer.Port)
	if err != nil {
		return nil, err
	}
	s := &Server{
		notify:   errChan,
		cfg:      cfg,
		logger:   logger,
		task:     t,
		listener: ln,
	}
	s.start(ctx)
	return s, nil
}

func (s *Server) start(ctx context.Context) {
	go func() {
		s.notify <- s.socketStart(ctx)

	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.listener.Close()
}
