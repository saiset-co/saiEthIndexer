package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase"
	"go.uber.org/zap"
)

// NewRouter -.
// Swagger spec:
// @title       Go boilerplate microservice framework
// @description Go boilerplate microservice framework
// @version     1.0
// @host        localhost:8081 (todo:dynamic host)
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l *zap.Logger, u *usecase.SomeUseCase) {
	handler.Use(GinLogger(l), GinRecovery(l, false), AuthRequired(l))

	ucHandler := &someHandler{
		uc:     u,
		logger: l,
	}

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Routers
	h := handler.Group("/v1")

	{
		h.GET("/get", ucHandler.get)
		h.POST("/post", ucHandler.set)

	}
}
