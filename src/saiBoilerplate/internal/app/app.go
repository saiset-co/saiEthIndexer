package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiBoilerplate/config"
	v1 "github.com/webmakom-com/saiBoilerplate/internal/handlers/http/v1"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase/repo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func Run(cfg *config.Configuration) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error when start logger : %s", err)
	}

	// mongo db repository
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// todo:password ?
	mongoClientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.Storage.Host, cfg.Storage.Port))

	client, err := mongo.Connect(ctx, mongoClientOptions)
	if err != nil {
		logger.Fatal("error when connect to mongo :", zap.Error(err))
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal("error when ping mongo instance :", zap.Error(err))
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			logger.Fatal("error when disconnect to mongo instance :", zap.Error(err))
		}
	}()

	mongoCollection := client.Database(cfg.Storage.Database).Collection(cfg.Storage.Collection)

	someUseCase := usecase.New(
		repo.New(mongoCollection),
	)
	handler := gin.New()
	v1.NewRouter(handler, logger, someUseCase)

}
