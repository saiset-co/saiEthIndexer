package handlers

import (
	"fmt"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/tasks"
	"go.uber.org/zap"
)

type HttpHandler struct {
	Logger      *zap.Logger
	TaskManager *tasks.TaskManager
}

type addContractsRequest struct {
	Contracts []config.Contract `json:"contracts" valid:",required"`
}

// Validation of contracts struct
func (r *addContractsRequest) validate() error {
	_, err := valid.ValidateStruct(r)

	return err
}

type addContractResponse struct {
	Created bool `json:"is_added" example:"true"`
}

func HandleHTTP(g *gin.RouterGroup, logger *zap.Logger, t *tasks.TaskManager) {
	handler := &HttpHandler{
		Logger:      logger,
		TaskManager: t,
	}
	{
		g.POST("/add_contract", handler.addContract)
	}
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
// @Router      /add_contract [post]
func (h *HttpHandler) addContract(c *gin.Context) {
	dto := addContractsRequest{}
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		h.Logger.Error("http  - add contract - bind", zap.Error(err))
		c.JSON(http.StatusBadRequest, errBadRequest)
	}
	fmt.Println(dto)

	for _, contract := range dto.Contracts {
		err = contract.Validate()
		if err != nil {
			h.Logger.Error("http  - add contract - validate", zap.Error(err))
			c.JSON(http.StatusBadRequest, errBadRequest)
			return
		}
	}
	err = h.TaskManager.AddContract(dto.Contracts)
	if err != nil {
		h.Logger.Error("http - v1 - set - repo", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errInternalServer)
		return
	}

	c.JSON(http.StatusOK, &addContractResponse{Created: true})
}
