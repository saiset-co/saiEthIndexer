package v1

import (
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiBoilerplate/internal/entity"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase"
	"go.uber.org/zap"
)

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

	c.JSON(http.StatusOK, someResponse{somes})
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
// @Router      /set [set]
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
