package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase"
	"go.uber.org/zap"
)

func NewRouter(handler *gin.Engine, l *zap.Logger, u *usecase.SomeUseCase) {

	ucHandler := &someHandler{
		uc:     u,
		logger: l,
	}

	// Routers
	h := handler.Group("/v1")

	{
		h.GET("/get", ucHandler.get)
		h.POST("/post", ucHandler.set)
		h.GET("/ws", ucHandler.websocket)

	}
}
