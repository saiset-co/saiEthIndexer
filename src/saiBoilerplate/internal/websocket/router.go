package websocket

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/webmakom-com/saiBoilerplate/handlers"
	"github.com/webmakom-com/saiBoilerplate/internal/http"
	"github.com/webmakom-com/saiBoilerplate/tasks"
	"go.uber.org/zap"
)

// NewRouter
// Swagger spec:
// @title       Go boilerplate microservice framework
// @description Go boilerplate microservice framework
// @version     1.0
// @host        localhost:8081
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l *zap.Logger, t *tasks.Task) {
	handler.Use(http.GinLogger(l), http.GinRecovery(l, false), http.AuthRequired(l)) // middlewares from http package

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Routers

	handlers.HandleWebsocket(handler, t, l)
}
