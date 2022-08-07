package websocket

import (
	"context"
	"net/http"
	"time"

	"github.com/webmakom-com/saiBoilerplate/config"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = ":80"
	defaultShutdownTimeout = 3 * time.Second
)

// Server struct
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// New returns new instance of http server for websocket
func New(handler http.Handler, cfg *config.Configuration, errChan chan error) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         cfg.Common.WebSocket.Url,
	}

	s := &Server{
		server:          httpServer,
		notify:          errChan,
		shutdownTimeout: defaultShutdownTimeout,
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
	}()
}

// Notify
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
