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
	"github.com/webmakom-com/saiP2P/config"
	"github.com/webmakom-com/saiP2P/handlers"
	"github.com/webmakom-com/saiP2P/internal/http"
	"github.com/webmakom-com/saiP2P/internal/socket"
	"github.com/webmakom-com/saiP2P/internal/websocket"
	"github.com/webmakom-com/saiP2P/tasks"
	"go.uber.org/zap"
)

type App struct {
	Cfg         *config.Configuration
	Logger      *zap.Logger
	taskManager *tasks.TaskManager
	handlers    *handlers.Handlers
}

func New() *App {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		log.Fatalf("error when start logger : %s", err)
	}
	return &App{
		Logger: logger,
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
	fmt.Printf("%+v\n", a.Cfg) // debug
	return nil
}

// Register task to app (main business logic)
func (a *App) RegisterTask(task *tasks.TaskManager) {
	a.taskManager = task
	return
}

// Register handlers to app
func (a *App) RegisterHandlers() {
	multihandler := handlers.Handlers{}
	if a.Cfg.Common.HttpServer.Enabled {
		//http server
		handler := gin.New()
		http.NewRouter(handler, a.Logger, a.taskManager)
		multihandler.Http = handler

	}

	if a.Cfg.Common.WebSocket.Enabled {
		// websocket server
		wsHandler := gin.New()
		websocket.NewRouter(wsHandler, a.Logger, a.taskManager)
		multihandler.Websocket = wsHandler
	}

	a.handlers = &multihandler
}

func (a *App) Run() error {
	errChan := make(chan error, 1)
	var (
		socketServer    = &socket.Server{}
		httpServer      = &http.HttpServer{}
		websocketServer = &websocket.Server{}
		err             error
	)
	if a.Cfg.Common.SocketServer.Enabled {
		socketServer, err = socket.New(context.Background(), a.Cfg, a.Logger, a.taskManager, errChan)
		if err != nil {
			return err
		}

	}

	if a.Cfg.Common.HttpServer.Enabled {
		httpServer = http.New(a.handlers.Http, a.Cfg, errChan)
	}

	if a.Cfg.Common.WebSocket.Enabled {
		websocketServer = websocket.New(a.handlers.Websocket, a.Cfg, errChan)
	}

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		a.Logger.Error("app - Run - signal: " + s.String())
	case err := <-errChan:
		a.Logger.Error("app - Run - server notifier: ", zap.Error(err))
	}
	if a.Cfg.Common.SocketServer.Enabled {
		err := socketServer.Shutdown()
		if err != nil {
			a.Logger.Error("app - Run - socketServer.Shutdown: ", zap.Error(err))
			return err
		}
		a.Logger.Info("socket server shutdown")
	}
	if a.Cfg.Common.HttpServer.Enabled {
		err := httpServer.Shutdown()
		if err != nil {
			a.Logger.Error("app - Run - httpServer.Shutdown: ", zap.Error(err))
			return err
		}
		a.Logger.Info("http server shutdown")
	}

	if a.Cfg.Common.WebSocket.Enabled {
		err = websocketServer.Shutdown()
		if err != nil {
			a.Logger.Error("app - Run - websocketServer.Shutdown: ", zap.Error(err))
		}
		a.Logger.Info("websocket server shutdown")
	}

	return nil
}
