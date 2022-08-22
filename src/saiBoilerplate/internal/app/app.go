package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/webmakom-com/saiBoilerplate/config"
	"github.com/webmakom-com/saiBoilerplate/handlers"
	"github.com/webmakom-com/saiBoilerplate/internal/http"
	"github.com/webmakom-com/saiBoilerplate/pkg/httpserver"
	"github.com/webmakom-com/saiBoilerplate/tasks"
	"go.uber.org/zap"
)

type App struct {
	Cfg         *config.Configuration
	logger      *zap.Logger
	handlers    *handlers.Handlers
	taskManager *tasks.TaskManager
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

	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("config read error: %w", err)
	}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return fmt.Errorf("config unmarshal error: %w", err)
	}

	a.Cfg = &cfg
	fmt.Printf("start config :%+v\n", a.Cfg) // debug
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
		http.NewRouter(handler, a.logger, a.taskManager)
		multihandler.Http = handler

	}

	a.handlers = &multihandler
}

func (a *App) Run() error {
	errChan := make(chan error, 1)
	var (
		httpServer = &httpserver.Server{}
	)
	if a.Cfg.Common.HttpServer.Enabled {
		httpServer = httpserver.New(a.handlers.Http, a.Cfg)
	}

	go a.taskManager.ProcessBlocks()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		a.logger.Error("app - Run - signal: " + s.String())
	case err := <-errChan:
		a.logger.Error("app - Run - server notifier: ", zap.Error(err))
	}
	if a.Cfg.Common.HttpServer.Enabled {
		err := httpServer.Shutdown()
		if err != nil {
			a.logger.Error("app - Run - httpServer.Shutdown: ", zap.Error(err))
			return err
		}
		a.logger.Info("http server shutdown")
	}

	return nil
}
