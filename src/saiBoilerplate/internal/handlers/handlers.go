package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiBoilerplate/internal/handlers/socket"
)

type Handler struct {
	SocketHandler *socket.Handler
	HttpHandler   *gin.Engine
	WsHandler     *gin.Engine
}
