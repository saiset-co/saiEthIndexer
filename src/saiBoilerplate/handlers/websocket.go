package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/tasks"
	"github.com/webmakom-com/saiBoilerplate/types"
	"go.uber.org/zap"
)

type WebsocketHandler struct {
	task   *tasks.Task
	logger *zap.Logger
	cfg    *config.Configuration
}

func HandleWebsocket(g *gin.Engine, t *tasks.Task, logger *zap.Logger) {
	handler := &WebsocketHandler{
		logger: logger,
		task:   t,
	}
	{
		g.GET("/", handler.handle)
	}
}

// @Summary     Simple Get and Set through websocket
// @Description Simple Get and Set through websocket
// @ID          Simple Get and Set through websocket
// @Tags  	    some
// @Accept      json
// @Produce     json
// @Success     200 {object} someResponse
// @Failure     500 {object} errInternalServerErr
func (h *WebsocketHandler) handle(c *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// token := r.Header.Get("token")
			// if token == "" {
			// 	return false
			// }
			// if token != h.cfg.Common.WebSocket.Token {
			// 	return false
			// }
			return true
		},
	}

	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("websocket - upgrade connection", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errInternalServer)
		return
	}

	for {
		msgType, b, err := connection.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			h.logger.Error("websocket - read message", zap.Error(err))
			break
		}

		var msg httpMessage
		buf := bytes.NewBuffer(b)
		err = json.Unmarshal(buf.Bytes(), &msg)
		if err != nil {
			h.logger.Error("websocket - Unmarshal", zap.Error(err))
			continue
		}
		switch msg.Method {
		case getMethod:
			somes, err := h.task.GetAll(c.Request.Context())
			if err != nil {
				h.logger.Error("websocket - get", zap.Error(err))
				continue
			}
			respBytes, err := json.Marshal(somes)
			if err != nil {
				h.logger.Error("websocket - marshal somes", zap.Error(err))
				continue
			}
			err = connection.WriteMessage(websocket.TextMessage, respBytes)
			if err != nil {
				h.logger.Error("websocket - write get answer", zap.Error(err))
				continue
			}
		case setMethod:
			some := types.Some{
				Key: msg.Key,
			}
			err := h.task.Set(c.Request.Context(), &some)
			if err != nil {
				h.logger.Error("socket - socketStart - set", zap.Error(err))
				continue
			}
			err = connection.WriteMessage(websocket.TextMessage, []byte("ok"))
			if err != nil {
				h.logger.Error("websocket - write set answer", zap.Error(err))
				continue
			}
		default:
			h.logger.Error("websocket - unknown method", zap.Error(errors.New("Unknown method : "+msg.Method)))
			err = connection.WriteMessage(websocket.TextMessage, []byte("unknown method : "+msg.Method))
			if err != nil {
				h.logger.Error("websocket - unknown method - write set answer", zap.Error(err))
				continue
			}
			continue
		}
	}
}
