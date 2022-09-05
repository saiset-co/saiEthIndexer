package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/webmakom-com/saiP2P/handlers"
	"github.com/webmakom-com/saiP2P/tasks"
	"go.uber.org/zap"
)

// NewRouter
// Swagger spec:
// @title       Go boilerplate microservice framework
// @description Go boilerplate microservice framework
// @version     1.0
// @host        localhost:8081
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l *zap.Logger, t *tasks.TaskManager) {
	handler.Use(GinLogger(l), GinRecovery(l, false), AuthRequired(l))

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	g := handler.Group("/v1")

	// func to realize in handlers package
	handlers.HandleHTTP(g, t, l)
}
