// Package httpserver implements socket server.
package socketserver

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"

	"github.com/webmakom-com/saiBoilerplate/config"
	"go.uber.org/zap"
)

const (
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	notify          chan error
	shutdownTimeout time.Duration
	logger          *zap.Logger
	cfg             *config.Configuration
	BufChannel      chan []byte
}

func New(cfg *config.Configuration, logger *zap.Logger) *Server {
	s := &Server{
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
		cfg:             cfg,
		logger:          logger,
		BufChannel:      make(chan []byte, 100000),
	}
	s.start()
	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.socketStart()
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) socketStart() error {
	ln, err := net.Listen("tcp", s.cfg.SocketServer.Host+":"+s.cfg.SocketServer.Port)
	if err != nil {
		return err
	}
	defer ln.Close()
newConn:
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					goto newConn
				}
				s.logger.Error("socket - start - accept", zap.Error(err))
				continue
			}
			s.logger.Info("socket - start - message", zap.String("message", message))
			_, err = conn.Write([]byte("got message : " + message))
			if err != nil {
				s.logger.Error("socket - start - write", zap.Error(err))
				continue
			}
			continue
		}

	}
}
