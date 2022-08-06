package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"net"

	"github.com/webmakom-com/saiBoilerplate/tasks"
	"github.com/webmakom-com/saiBoilerplate/types"
	"go.uber.org/zap"
)

const (
	getMethod = "get"
	setMethod = "set"
)

type socketHandler struct {
	logger     *zap.Logger
	connection net.Conn
	task       *tasks.Task
}

// socket message
type socketMessage struct {
	Method string `json:"method"`
	Token  string `json:"token"`
	Key    string `json:"key"`
}

// socket handling func. Should with signature func (context.Context, b []byte)
func HandleSocket(ctx context.Context, conn net.Conn, b []byte, logger *zap.Logger, t *tasks.Task) error {
	s := &socketHandler{
		logger:     logger,
		task:       t,
		connection: conn,
	}
	var msg socketMessage
	buf := bytes.NewBuffer(b)
	err := json.Unmarshal(buf.Bytes(), &msg)
	if err != nil {
		s.logger.Error("socket - socketStart - Unmarshal", zap.Error(err))
		return err

	}
	//dumb auth check
	if msg.Token == "" {
		s.logger.Error("socket - socketStart - auth", zap.Error(errors.New("auth failed:empty token")))
		return err
	}
	switch msg.Method {
	case getMethod:
		somes, err := s.task.GetAll(ctx)
		if err != nil {
			s.logger.Error("socket - socketStart - get", zap.Error(err))
			return err
		}
		respBytes, err := json.Marshal(somes)
		if err != nil {
			s.logger.Error("socket - socketStart - marshal somes", zap.Error(err))
			return err
		}
		_, err = s.connection.Write(respBytes)
		if err != nil {
			s.logger.Error("socket - socketStart - write get answer", zap.Error(err))
			return err
		}
	case setMethod:
		some := types.Some{
			Key: msg.Key,
		}
		err := s.task.Set(ctx, &some)
		if err != nil {
			s.logger.Error("socket - socketStart - set", zap.Error(err))
			return err
		}
		_, err = s.connection.Write([]byte("ok"))
		if err != nil {
			s.logger.Error("socket - socketStart - write set answer", zap.Error(err))
			return err
		}
	default:
		s.logger.Error("socket - socketStart - unknown method", zap.Error(errors.New("Unknown method : "+msg.Method)))
		_, err = s.connection.Write([]byte("unknown method : " + msg.Method))
		if err != nil {
			s.logger.Error("socket - socketStart - unknown method - write set answer", zap.Error(err))
			return err
		}
		return err
	}
	return nil
}
