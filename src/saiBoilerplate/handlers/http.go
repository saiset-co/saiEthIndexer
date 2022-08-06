package handlers

import (
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiBoilerplate/tasks"
	"github.com/webmakom-com/saiBoilerplate/types"
	"go.uber.org/zap"
)

type httpMessage struct {
	Method string `json:"method"`
	Token  string `json:"token"`
	Key    string `json:"key"`
}

type HttpHandler struct {
	Logger *zap.Logger
	Task   *tasks.Task
}

type someResponse struct {
	Somes []*types.Some `json:"Somes"`
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

func HandleHTTP(g *gin.RouterGroup, t *tasks.Task, logger *zap.Logger) {
	handler := &HttpHandler{
		Logger: logger,
		Task:   t,
	}
	{
		g.GET("/get", handler.get)
		g.POST("/post", handler.set)
	}
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
func (h *HttpHandler) get(c *gin.Context) {
	somes, err := h.Task.GetAll(c.Request.Context())
	if err != nil {
		h.Logger.Error("http - v1 - get", zap.Error(err))
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
func (h *HttpHandler) set(c *gin.Context) {
	dto := &setRequest{}
	err := c.ShouldBindJSON(dto)
	if err != nil {
		h.Logger.Error("http - v1 - set - bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, errBadRequest)
	}
	some := &types.Some{
		Key: dto.Key,
	}

	err = h.Task.Set(c.Request.Context(), some)
	if err != nil {
		h.Logger.Error("http - v1 - set - repo", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errInternalServer)
		return
	}

	c.JSON(http.StatusOK, &setResponse{Created: true})
}
