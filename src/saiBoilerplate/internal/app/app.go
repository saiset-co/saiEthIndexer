package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/internal/handlers"
	v1 "github.com/webmakom-com/saiBoilerplate/internal/handlers/http/v1"
	"github.com/webmakom-com/saiBoilerplate/internal/handlers/socket"
	"github.com/webmakom-com/saiBoilerplate/internal/handlers/websocket"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase"
	"github.com/webmakom-com/saiBoilerplate/internal/usecase/repo"
	"github.com/webmakom-com/saiBoilerplate/pkg/httpserver"
	websocketserver "github.com/webmakom-com/saiBoilerplate/pkg/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type App struct {
	cfg      *config.Configuration
	logger   *zap.Logger
	repo     *repo.SomeRepo
	uc       *usecase.SomeUseCase
	handlers *handlers.Handler
}

func New() *App {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("error when start logger : %s", err)
	}
	return &App{
		logger: logger,
	}
}

// Register config to app
func (a *App) RegisterConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	fmt.Printf("loaded configuration:%+v\n", cfg)

	a.cfg = &cfg

	return nil
}

// Register storage to app
func (a *App) RegisterStorage() error {
	ctx := context.Background()

	// use mongodb as a storage
	mongoClientOptions := &options.ClientOptions{}
	if a.cfg.Mongo.User != "" && a.cfg.Mongo.Pass != "" {
		mongoClientOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", a.cfg.Mongo.Host, a.cfg.Mongo.Port)).SetAuth(options.Credential{
			Username: a.cfg.Mongo.User,
			Password: a.cfg.Mongo.Pass,
		})
	} else {
		mongoClientOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", a.cfg.Mongo.Host, a.cfg.Mongo.Port))
	}
	client, err := mongo.Connect(ctx, mongoClientOptions)
	if err != nil {
		a.logger.Fatal("error when connect to mongo :", zap.Error(err))
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		a.logger.Fatal("error when ping mongo instance :", zap.Error(err))
		return err
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			a.logger.Fatal("error when disconnect to mongo instance :", zap.Error(err))
		}
	}()
	mongoCollection := client.Database(a.cfg.Mongo.Database).Collection(a.cfg.Mongo.Collection)

	a.logger.Info("found collection", zap.String("mongo collection", mongoCollection.Name()))

	repo := repo.New(mongoCollection)

	a.repo = repo

	return nil

}

// Register usecase (tasks) to app (main business logic)
func (a *App) RegisterUsecase() {
	someUseCase := usecase.New(a.repo)
	a.uc = someUseCase
}

// Register handlers to app
func (a *App) RegisterHandlers() {
	if a.cfg.Common.SocketServer.Enabled {
		socketHandler := socket.New(context.Background(), a.cfg, a.logger, a.uc)
		a.handlers.SocketHandler = socketHandler
	}

	if a.cfg.Common.HttpServer.Enabled {
		//http server
		handler := gin.New()
		v1.NewRouter(handler, a.logger, a.uc)
		a.handlers.HttpHandler = handler

	}

	if a.cfg.Common.WebSocket.Enabled {
		// websocket server
		wsHandler := gin.New()
		websocket.NewRouter(wsHandler, a.logger, a.uc, a.cfg)
		a.handlers.WsHandler = wsHandler
	}
}

func (a *App) Run() {

	httpServer := httpserver.New(a.handlers.HttpHandler, a.cfg)

	websocketServer := websocketserver.New(a.handlers.WsHandler, a.cfg)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		a.logger.Error("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		a.logger.Error("app - Run - httpServer.Notify: ", zap.Error(err))
	case err := <-a.handlers.SocketHandler.Notify():
		a.logger.Error("app - Run - socketServer.Notify: ", zap.Error(err))
	case err := <-websocketServer.Notify():
		a.logger.Error("app - Run - websocketServer.Notify: ", zap.Error(err))

	}
	err := httpServer.Shutdown()
	if err != nil {
		a.logger.Error("app - Run - httpServer.Shutdown: ", zap.Error(err))
	}
}
