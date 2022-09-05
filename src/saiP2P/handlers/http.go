package handlers

import (
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiP2P/tasks"
	"github.com/webmakom-com/saiP2P/types"
	"go.uber.org/zap"
)

const (
	httpSource = "http"
)

type httpMessage struct {
	Method string `json:"method"`
	Token  string `json:"token"`
	Key    string `json:"key"`
}

type HttpHandler struct {
	Logger      *zap.Logger
	TaskManager *tasks.TaskManager
}

type msgResponse struct {
	Sent bool `json:"sent"`
}

type msg struct {
	Info string `json:"info" valid:",required"`
}

// Validation of incoming struct
func (r *msg) validate() error {
	_, err := valid.ValidateStruct(r)

	return err
}

type setResponse struct {
	Created bool `json:"created" example:"true"`
}

func HandleHTTP(g *gin.RouterGroup, t *tasks.TaskManager, logger *zap.Logger) {
	handler := &HttpHandler{
		Logger:      logger,
		TaskManager: t,
	}
	{
		g.POST("/msg", handler.msg)
	}
}

// @Summary     handle message from another node and send to callbacks
// @Description handle message from another node and send to callbacks
// @ID          handle message from another node and send to callbacks
// @Tags  	    CallbackMessage
// @Accept      json
// @Produce     json
// @Success     200 {object} Sent
// @Failure     500 {object} errInternalServerErr
// @Router      /msg [post]
func (h *HttpHandler) msg(c *gin.Context) {
	dto := &msg{}
	err := c.ShouldBindJSON(dto)
	if err != nil {
		h.Logger.Error("http - v1 - msg - bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, errBadRequest)
	}

	callbackMsg := &types.CallbackMessage{
		Source:  httpSource,
		Message: dto.Info,
	}
	h.TaskManager.SendCallbackMsg(callbackMsg)

	c.IndentedJSON(http.StatusOK, &msgResponse{Sent: true})
}
