package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/webmakom-com/saiBoilerplate/internal/entity"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase"
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

type someHandler struct {
	uc     *usecase.SomeUseCase
	logger *zap.Logger
}

type someResponse struct {
	Somes []*entity.Some `json:"Somes"`
}

type setRequest struct {
	Key string `json:"key" valid:",required"`
}

// Validation of incoming struct
func (r *setRequest) validate() error {
	_, err := valid.ValidateStruct(r)

	return err
}

type setResponse struct {
	Created bool `json:"created" example:"true"`
}

// @Summary     Simple Get
// @Description Simple Get
// @ID          Simple Get
// @Tags  	    some
// @Accept      json
// @Produce     json
// @Success     200 {object} someResponse
// @Failure     500 {object} errInternalServerErr
// @Router      /get [get]
func (h *someHandler) get(c *gin.Context) {
	somes, err := h.uc.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error("http - v1 - get", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errInternalServer)
		return
	}

	c.IndentedJSON(http.StatusOK, someResponse{somes})
}

// @Summary     Simple set
// @Description Simple set
// @ID          Simple set
// @Tags  	    some
// @Accept      json
// @Produce     json
// @Success     200 {object} setResponse
// @Failure     500 {object} errInternalServer
// @Failure     400 {object} errBadRequest
// @Router      /set [post]
func (h *someHandler) set(c *gin.Context) {
	dto := &setRequest{}
	err := c.ShouldBindJSON(dto)
	if err != nil {
		h.logger.Error("http - v1 - set - bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, errBadRequest)
	}
	some := &entity.Some{
		Key: dto.Key,
	}

	err = h.uc.Set(c.Request.Context(), some)
	if err != nil {
		h.logger.Error("http - v1 - set - repo", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errInternalServer)
		return
	}

	c.JSON(http.StatusOK, &setResponse{Created: true})
}

// @Summary     Simple websocket
// @Description Simple websocket
// @ID          Simple websocket
// @Tags  	    some
// @Accept      json
// @Produce     json
// @Success     200 {object} setResponse
// @Failure     500 {object} errInternalServer
// @Failure     400 {object} errBadRequest
// @Router      /ws [get]
func (h *someHandler) websocket(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // dumb auth check
		},
	}
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("http - v1 - ws - upgrade", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errInternalServer)
		return
	}

	defer wsConn.Close()

	for {
		//todo:check message type
		mt, b, err := wsConn.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			h.logger.Error("http - v1 - ws - upgrade - read message", zap.Error(err))
			break
		}
		var msg socketMessage
		buf := bytes.NewBuffer(b)
		err = json.Unmarshal(buf.Bytes(), &msg)
		if err != nil {
			h.logger.Error("http - v1 - ws - upgrade - read message - Unmarshal", zap.Error(err))
			continue
		}
		//dumb auth check
		if msg.Token == "" {
			h.logger.Error("http - v1 - ws - upgrade - read message - Unmarshal - auth", zap.Error(errors.New("auth failed:empty token")))
			continue
		}
		switch msg.Method {
		case getMethod:
			somes, err := h.uc.GetAll(c.Request.Context())
			if err != nil {
				h.logger.Error("http - v1 - ws - upgrade - read message - Unmarshal - get", zap.Error(err))
				continue
			}
			respBytes, err := json.Marshal(somes)
			if err != nil {
				h.logger.Error("shttp - v1 - ws - upgrade - read message - Unmarshal - marshal somes", zap.Error(err))
				continue
			}
			err = wsConn.WriteMessage(websocket.TextMessage, respBytes)
			if err != nil {
				h.logger.Error("http - v1 - ws - upgrade - read message - Unmarshal - write get answer", zap.Error(err))
				continue
			}
		case setMethod:
			some := entity.Some{
				Key: msg.Key,
			}
			err := h.uc.Set(c.Request.Context(), &some)
			if err != nil {
				h.logger.Error("http - v1 - ws - upgrade - read message - Unmarshal - set", zap.Error(err))
				continue
			}
			err = wsConn.WriteMessage(websocket.TextMessage, []byte("ok"))
			if err != nil {
				h.logger.Error("http - v1 - ws - upgrade - read message - Unmarshal - write set answer", zap.Error(err))
				continue
			}
		default:
			h.logger.Error("http - v1 - ws - upgrade - read message - Unmarshal - unknown method", zap.Error(errors.New("Unknown method : "+msg.Method)))
			err = wsConn.WriteMessage(websocket.TextMessage, []byte("unknown method : "+msg.Method))
			if err != nil {
				h.logger.Error("http - v1 - ws - upgrade - read message - Unmarshal - unknown method - write set answer", zap.Error(err))
				continue
			}
			continue

		}

	}
}
