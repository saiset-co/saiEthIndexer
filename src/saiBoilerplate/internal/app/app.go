package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiBoilerplate/config"
	v1 "github.com/webmakom-com/saiBoilerplate/internal/handlers/http/v1"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase/repo"
	"github.com/webmakom-com/saiBoilerplate/pkg/httpserver"
	"github.com/webmakom-com/saiBoilerplate/pkg/socketserver"
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

	mongoClientOptions := &options.ClientOptions{}
	if cfg.Mongo.User != "" && cfg.Mongo.Pass != "" {
		mongoClientOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.Mongo.Host, cfg.Mongo.Port)).SetAuth(options.Credential{
			Username: cfg.Mongo.User,
			Password: cfg.Mongo.Pass,
		})
	} else {
		mongoClientOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", cfg.Mongo.Host, cfg.Mongo.Port))
	}
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
	mongoCollection := client.Database(cfg.Mongo.Database).Collection(cfg.Mongo.Collection)

	logger.Info("found collection", zap.String("mongo collection", mongoCollection.Name()))

	someUseCase := usecase.New(
		repo.New(mongoCollection),
	)

	//http server
	handler := gin.New()
	v1.NewRouter(handler, logger, someUseCase)

	httpServer := httpserver.New(handler, cfg)

	// socket server

	socketServer := socketserver.New(cfg, logger)
	//go socket.Handle(socketServer.BufChannel)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Error("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		logger.Error("app - Run - httpServer.Notify: ", zap.Error(err))
	case err = <-socketServer.Notify():
		logger.Error("app - Run - socketServer.Notify: ", zap.Error(err))

	}
	err = httpServer.Shutdown()
	if err != nil {
		logger.Error("app - Run - httpServer.Shutdown: ", zap.Error(err))
	}
}
