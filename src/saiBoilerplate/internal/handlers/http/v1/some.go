package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase"
	"go.uber.org/zap"
)

type someHandler struct {
	uc     usecase.SomeUseCase
	logger *zap.Logger
}

// @Summary     Simple Get
// @Description Simple Get
// @ID          Simple Get
// @Tags  	    some
// @Accept      json
// @Produce     json
// @Success     200 {object} historyResponse
// @Failure     500 {object} response
// @Router      /get [get]
func (h *someHandler) get(c *gin.Context) {
	somes, err := h.uc.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error("http - v1 - get", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errInternalServerErr)
		return
	}

	c.JSON(http.StatusOK, somes)
}
