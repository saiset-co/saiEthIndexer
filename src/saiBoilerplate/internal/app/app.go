package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/handlers"
	"github.com/webmakom-com/saiBoilerplate/internal/http"
	"github.com/webmakom-com/saiBoilerplate/internal/socket"
	"github.com/webmakom-com/saiBoilerplate/internal/websocket"
	"github.com/webmakom-com/saiBoilerplate/storage"
	"github.com/webmakom-com/saiBoilerplate/tasks"
	"github.com/webmakom-com/saiBoilerplate/tasks/repo"
	"go.uber.org/zap"
)

type App struct {
	Cfg      *config.Configuration
	logger   *zap.Logger
	task     *tasks.Task
	repo     *repo.SomeRepo
	handlers *handlers.Handlers
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
func (a *App) RegisterConfig(path string) error {
	cfg := config.Configuration{}

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	a.Cfg = &cfg

	return nil
}

// Register storage to app
func (a *App) RegisterStorage(storage *storage.Storage) error {
	a.repo = &repo.SomeRepo{
		Collection: storage.Collection,
	}
	return nil

}

// Register task to app (main business logic)
func (a *App) RegisterTask(task *tasks.Task) {
	a.task = task
}

// Register handlers to app
func (a *App) RegisterHandlers() {

	if a.Cfg.Common.HttpServer.Enabled {
		//http server
		handler := gin.New()
		http.NewRouter(handler, a.logger, a.task)
		a.handlers.Http = handler

	}

	if a.Cfg.Common.WebSocket.Enabled {
		// websocket server
		wsHandler := gin.New()
		websocket.NewRouter(wsHandler, a.logger, a.task)
		a.handlers.Websocket = wsHandler
	}
}

func (a *App) Run() {
	var socketHandler = &socket.Handler{}
	if a.Cfg.Common.SocketServer.Enabled {
		socketHandler = socket.New(context.Background(), a.Cfg, a.logger, a.task)
	}
	httpServer := http.New(a.handlers.Http, a.Cfg)

	websocketServer := websocket.New(a.handlers.Websocket, a.Cfg)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		a.logger.Error("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		a.logger.Error("app - Run - httpServer.Notify: ", zap.Error(err))
	case err := <-socketHandler.Notify():
		a.logger.Error("app - Run - socketServer.Notify: ", zap.Error(err))
	case err := <-websocketServer.Notify():
		a.logger.Error("app - Run - websocketServer.Notify: ", zap.Error(err))

	}
	err := httpServer.Shutdown()
	if err != nil {
		a.logger.Error("app - Run - httpServer.Shutdown: ", zap.Error(err))
	}
}
