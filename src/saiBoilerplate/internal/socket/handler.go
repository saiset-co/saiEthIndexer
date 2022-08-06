package socket

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"

	"github.com/webmakom-com/saiBoilerplate/handlers"
	"go.uber.org/zap"
)

const (
	getMethod = "get"
	setMethod = "set"
)

func (s *Handler) socketStart(ctx context.Context) error {

	ln, err := net.Listen("tcp", s.cfg.Common.SocketServer.Host+":"+s.cfg.Common.SocketServer.Port)
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

			//todo :insert socket handler here (simple func, mot method ???)

			err = handlers.HandleSocket(ctx, conn, b, s.logger, s.task)
			if err != nil {
				s.logger.Info("socket - handle", zap.Error(err))
				continue
			}

		}
	}
}
