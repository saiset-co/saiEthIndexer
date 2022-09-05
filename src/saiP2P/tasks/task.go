package tasks

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/webmakom-com/saiP2P/config"
	"github.com/webmakom-com/saiP2P/types"
	"go.uber.org/zap"
)

// Example struct
type TaskManager struct {
	Cfg    *config.Configuration
	Logger *zap.Logger
}

// New creates new taskManager instance
func New(config *config.Configuration, logger *zap.Logger) *TaskManager {
	return &TaskManager{
		Cfg:    config,
		Logger: logger,
	}
}

// Send callback messages to enabled callbacks in configuration file
func (t *TaskManager) SendCallbackMsg(msg *types.CallbackMessage) {
	if t.Cfg.Specific.HttpCallback.Enabled {
		err := t.SendHTTPMsg(msg)
		if err != nil {
			t.Logger.Error("taskManager - SendCallbackMsg - SendHTTPMsg", zap.Error(err))
		}
	}
	if t.Cfg.Specific.SocketCallback.Enabled {
		err := t.SendSocketMsg(msg)
		if err != nil {
			t.Logger.Error("taskManager - SendCallbackMsg - SendSocketMsg", zap.Error(err))
		}
	}
	if t.Cfg.Specific.WebsocketCallback.Enabled {
		err := t.SendWebsocketMsg(msg)
		if err != nil {
			t.Logger.Error("taskManager - SendCallbackMsg - SendWebsocketMsg", zap.Error(err))
		}
	}
}

func (t *TaskManager) SendHTTPMsg(msg *types.CallbackMessage) error {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Logger.Error("taskManager - SendCallbackMsg - SendHTTPMsg - Marshal msg", zap.Error(err))
		return err
	}

	req, err := http.NewRequest("POST", t.Cfg.Specific.HttpCallback.Address, bytes.NewBuffer(data))
	if err != nil {
		t.Logger.Error("taskManager - SendCallbackMsg - SendHTTPMsg - create request", zap.Error(err))
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Logger.Error("taskManager - SendCallbackMsg - SendHTTPMsg - do request", zap.Error(err))
		return err
	}

	if resp.StatusCode < 200 && resp.StatusCode > 202 {
		t.Logger.Sugar().Infof("taskManager - SendCallbackMsg - SendHTTPMsg - wrong status code after send callback msg : %d", resp.StatusCode)
		return errors.New("wrong status code returned")
	}
	t.Logger.Sugar().Infof("Successfully sent callback message to http callback, msg : %s", msg)
	return nil
}

func (t *TaskManager) SendSocketMsg(msg *types.CallbackMessage) error {
	conn, err := net.Dial("tcp", t.Cfg.Specific.SocketCallback.Address)
	if err != nil {
		t.Logger.Error("taskManager - SendCallbackMsg - SendSocketMsg - dial", zap.Error(err))
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(msg)
	if err != nil {
		t.Logger.Error("taskManager - SendCallbackMsg - SendSocketMsg - Marshal msg", zap.Error(err))
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		t.Logger.Error("taskManager - SendCallbackMsg - SendSocketMsg - write msg", zap.Error(err))
		return err
	}
	t.Logger.Sugar().Infof("Successfully sent callback message to socket callback, msg : %s", msg)
	return nil
}

func (t *TaskManager) SendWebsocketMsg(msg *types.CallbackMessage) error {
	// todo : add auth header & handle response if needed
	conn, _, err := websocket.DefaultDialer.Dial(t.Cfg.Specific.WebsocketCallback.Address, nil)
	if err != nil {
		t.Logger.Error("taskManager - SendCallbackMsg - SendWebsocketMsg - dial websocket", zap.Error(err))
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(msg)
	if err != nil {
		t.Logger.Error("taskManager - SendCallbackMsg - SendWebsocketMsg - Marshal msg", zap.Error(err))
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		t.Logger.Error("taskManager - SendCallbackMsg - SendWebsocketMsg - write msg", zap.Error(err))
		return err
	}
	t.Logger.Sugar().Infof("Successfully sent callback message to websocket callback, msg : %s", msg)
	return nil
}
