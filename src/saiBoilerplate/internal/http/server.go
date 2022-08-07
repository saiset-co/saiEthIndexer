package http

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

// Server
type HttpServer struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// New returns new instance of http server
func New(handler http.Handler, cfg *config.Configuration, errChan chan error) *HttpServer {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         cfg.Common.HttpServer.Host + ":" + cfg.Common.HttpServer.Port,
	}

	s := &HttpServer{
		server:          httpServer,
		notify:          errChan,
		shutdownTimeout: defaultShutdownTimeout,
	}

	s.start()

	return s
}

func (s *HttpServer) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
	}()
}

// Notify
func (s *HttpServer) Notify() <-chan error {
	return s.notify
}

// Shutdown
func (s *HttpServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
