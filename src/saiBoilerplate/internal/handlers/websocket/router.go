package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase"
	"go.uber.org/zap"
)

func NewRouter(handler *gin.Engine, l *zap.Logger, u *usecase.SomeUseCase, cfg *config.Configuration) {

	wsHandler := &someWSHandler{
		uc:     u,
		logger: l,
		cfg:    cfg,
	}

	// Routers
	handler.GET(cfg.WebSocket.Url, wsHandler.handle)
}
