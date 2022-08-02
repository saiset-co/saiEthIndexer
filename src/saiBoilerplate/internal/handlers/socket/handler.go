package socket

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"

	"github.com/webmakom-com/saiBoilerplate/internal/entity"
	"go.uber.org/zap"
)

const (
	getMethod = "get"
	setMethod = "set"
)

type socketMessage struct {
	Method string `json:"method"`
	Token  string `json:"token"`
	Key    string `json:"key"`
}

func (s *Handler) socketStart(ctx context.Context) error {
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
			var msg socketMessage
			buf := bytes.NewBuffer(b)
			err = json.Unmarshal(buf.Bytes(), &msg)
			if err != nil {
				s.logger.Error("socket - socketStart - Unmarshal", zap.Error(err))
				continue
			}
			//dumb auth check
			if msg.Token == "" {
				s.logger.Error("socket - socketStart - auth", zap.Error(errors.New("auth failed:empty token")))
				continue
			}
			switch msg.Method {
			case getMethod:
				somes, err := s.uc.GetAll(ctx)
				if err != nil {
					s.logger.Error("socket - socketStart - get", zap.Error(err))
					continue
				}
				respBytes, err := json.Marshal(somes)
				if err != nil {
					s.logger.Error("socket - socketStart - marshal somes", zap.Error(err))
					continue
				}
				_, err = conn.Write(respBytes)
				if err != nil {
					s.logger.Error("socket - socketStart - write get answer", zap.Error(err))
					continue
				}
			case setMethod:
				some := entity.Some{
					Key: msg.Key,
				}
				err := s.uc.Set(ctx, &some)
				if err != nil {
					s.logger.Error("socket - socketStart - set", zap.Error(err))
					continue
				}
				_, err = conn.Write([]byte("ok"))
				if err != nil {
					s.logger.Error("socket - socketStart - write set answer", zap.Error(err))
					continue
				}
			default:
				s.logger.Error("socket - socketStart - unknown method", zap.Error(errors.New("Unknown method : "+msg.Method)))
				_, err = conn.Write([]byte("unknown method : " + msg.Method))
				if err != nil {
					s.logger.Error("socket - socketStart - unknown method - write set answer", zap.Error(err))
					continue
				}
				continue
			}

		}
	}
}
