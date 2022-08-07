package socket

import (
	"bufio"
	"context"
	"errors"
	"io"

	"github.com/webmakom-com/saiBoilerplate/handlers"
	"go.uber.org/zap"
)

const (
	getMethod = "get"
	setMethod = "set"
)

func (s *Server) socketStart(ctx context.Context) error {

	defer s.listener.Close()
newConn:
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()
		for {
			b, err := bufio.NewReader(conn).ReadBytes(byte('\n'))
			if err != nil {
				if errors.Is(err, io.EOF) {
					goto newConn
				}
				s.logger.Error("socket - start - accept", zap.Error(err))
				continue
			}
			s.logger.Info("socket - start - message", zap.String("message", string(b)))

			err = handlers.HandleSocket(ctx, conn, b, s.logger, s.task)
			if err != nil {
				s.logger.Info("socket - handle", zap.Error(err))
				continue
			}

		}
	}
}
