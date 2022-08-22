package tasks

import (
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (bm *BlockManager) SendWebSocketMsg(txMsg []byte, commandName string) error {
	header := make(http.Header)
	header.Add("Origin", bm.config.Common.HttpServer.Host)
	conn, _, err := websocket.DefaultDialer.Dial(bm.config.Common.WebSocket.Url, header)
	if err != nil {
		bm.logger.Error("tasks - SendWebSocketMg - dial with websocket server", zap.Error(err))
		return err
	}
	msg := []byte(bm.config.Common.WebSocket.Token + "|" + commandName + "-" + string(txMsg))
	err = conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		bm.logger.Error("tasks - SendWebSocketMg - write message to websocket server", zap.Error(err))
		return err
	}
	return nil
}
